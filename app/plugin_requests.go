// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"bytes"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/utils"
)

func (a *App) ServePluginRequest(w http.ResponseWriter, r *http.Request) {
	pluginsEnvironment := a.GetPluginsEnvironment()
	if pluginsEnvironment == nil {
		err := model.NewAppError("ServePluginRequest", "app.plugin.disabled.app_error", nil, "Enable plugins to serve plugin requests", http.StatusNotImplemented)
		a.Log.Error(err.Error())
		w.WriteHeader(err.StatusCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.ToJson()))
		return
	}

	params := mux.Vars(r)
	hooks, err := pluginsEnvironment.HooksForPlugin(params["plugin_id"])
	if err != nil {
		a.Log.Error("Access to route for non-existent plugin", mlog.String("missing_plugin_id", params["plugin_id"]), mlog.Err(err))
		http.NotFound(w, r)
		return
	}

	a.servePluginRequest(w, r, hooks.ServeHTTP)
}

// ServePluginStaticRequest serves static plugin files
// at the URL http(s)://$SITE_URL/plugins/$PLUGIN_ID/public/{anything}
func (a *App) ServePluginStaticRequest(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}

	publicPathStartIndex := strings.Index(r.URL.Path, "/public/")

	if publicPathStartIndex < 0 {
		http.NotFound(w, r)
		return
	}

	publicPathStartIndex += 8

	// Should be in the form of /$PLUGIN_ID/public/{anything} by the timne we get here
	pluginID := strings.Split(r.URL.Path, "/")[2]

	staticFiles, isOk := a.GetPluginsEnvironment().StaticFilesPath(pluginID)

	if !isOk {
		http.NotFound(w, r)
		return
	}

	requestedPublicFile := string([]rune(r.URL.Path)[publicPathStartIndex:])

	http.ServeFile(w, r, filepath.Join(staticFiles, requestedPublicFile))
}

func (a *App) servePluginRequest(w http.ResponseWriter, r *http.Request, handler func(*plugin.Context, http.ResponseWriter, *http.Request)) {
	token := ""
	context := &plugin.Context{
		RequestId:      model.NewId(),
		IpAddress:      utils.GetIpAddress(r),
		AcceptLanguage: r.Header.Get("Accept-Language"),
		UserAgent:      r.UserAgent(),
	}
	cookieAuth := false

	authHeader := r.Header.Get(model.HEADER_AUTH)
	if strings.HasPrefix(strings.ToUpper(authHeader), model.HEADER_BEARER+" ") {
		token = authHeader[len(model.HEADER_BEARER)+1:]
	} else if strings.HasPrefix(strings.ToLower(authHeader), model.HEADER_TOKEN+" ") {
		token = authHeader[len(model.HEADER_TOKEN)+1:]
	} else if cookie, _ := r.Cookie(model.SESSION_COOKIE_TOKEN); cookie != nil {
		token = cookie.Value
		cookieAuth = true
	} else {
		token = r.URL.Query().Get("access_token")
	}

	r.Header.Del("Mattermost-User-Id")
	if token != "" {
		session, err := a.GetSession(token)
		csrfCheckPassed := false

		if err == nil && cookieAuth && r.Method != "GET" {
			sentToken := ""

			if r.Header.Get(model.HEADER_CSRF_TOKEN) == "" {
				bodyBytes, _ := ioutil.ReadAll(r.Body)
				r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
				r.ParseForm()
				sentToken = r.FormValue("csrf")
				r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			} else {
				sentToken = r.Header.Get(model.HEADER_CSRF_TOKEN)
			}

			expectedToken := session.GetCSRF()

			if sentToken == expectedToken {
				csrfCheckPassed = true
			}

			// ToDo(DSchalla) 2019/01/04: Remove after deprecation period and only allow CSRF Header (MM-13657)
			if r.Header.Get(model.HEADER_REQUESTED_WITH) == model.HEADER_REQUESTED_WITH_XML && !csrfCheckPassed {
				csrfErrorMessage := "CSRF Check failed for request - Please migrate your plugin to either send a CSRF Header or Form Field, XMLHttpRequest is deprecated"
				if *a.Config().ServiceSettings.ExperimentalStrictCSRFEnforcement {
					a.Log.Warn(csrfErrorMessage)
				} else {
					a.Log.Debug(csrfErrorMessage)
					csrfCheckPassed = true
				}
			}
		} else {
			csrfCheckPassed = true
		}

		if session != nil && err == nil && csrfCheckPassed {
			r.Header.Set("Mattermost-User-Id", session.UserId)
			context.SessionId = session.Id
		}
	}

	cookies := r.Cookies()
	r.Header.Del("Cookie")
	for _, c := range cookies {
		if c.Name != model.SESSION_COOKIE_TOKEN {
			r.AddCookie(c)
		}
	}
	r.Header.Del(model.HEADER_AUTH)
	r.Header.Del("Referer")

	params := mux.Vars(r)

	subpath, _ := utils.GetSubpathFromConfig(a.Config())

	newQuery := r.URL.Query()
	newQuery.Del("access_token")
	r.URL.RawQuery = newQuery.Encode()
	r.URL.Path = strings.TrimPrefix(r.URL.Path, path.Join(subpath, "plugins", params["plugin_id"]))

	handler(context, w, r)
}
