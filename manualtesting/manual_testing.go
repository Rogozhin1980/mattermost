// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package manualtesting

import (
	l4g "code.google.com/p/log4go"
	"github.com/mattermost/platform/api"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
	"github.com/mattermost/platform/i18n"
	"hash/fnv"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type TestEnvironment struct {
	Params        map[string][]string
	Client        *model.Client
	CreatedTeamId string
	CreatedUserId string
	Context       *api.Context
	Writer        http.ResponseWriter
	Request       *http.Request
}

func InitManualTesting() {
	api.Srv.Router.Handle("/manualtest", api.AppHandler(manualTest)).Methods("GET")
}

func manualTest(c *api.Context, w http.ResponseWriter, r *http.Request) {
	T := i18n.GetSystemLanguage()
	// Let the world know
	l4g.Info("Setting up for manual test...")

	// URL Parameters
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		c.Err = model.NewAppError("/manual", "Unable to parse URL", "")
		return
	}

	// Grab a uuid (if available) to seed the random number generator so we don't get conflicts.
	uid, ok := params["uid"]
	if ok {
		hasher := fnv.New32a()
		hasher.Write([]byte(uid[0] + strconv.Itoa(int(time.Now().UTC().UnixNano()))))
		hash := hasher.Sum32()
		rand.Seed(int64(hash))
	} else {
		l4g.Debug("No uid in url")
	}

	// Create a client for tests to use
	client := model.NewClient("http://localhost" + utils.Cfg.ServiceSettings.ListenAddress)

	// Check for username parameter and create a user if present
	username, ok1 := params["username"]
	teamDisplayName, ok2 := params["teamname"]
	var teamID string
	var userID string
	if ok1 && ok2 {
		l4g.Info("Creating user and team")
		// Create team for testing
		team := &model.Team{
			DisplayName: teamDisplayName[0],
			Name:        utils.RandomName(utils.Range{20, 20}, utils.LOWERCASE),
			Email:       utils.RandomEmail(utils.Range{20, 20}, utils.LOWERCASE),
			Type:        model.TEAM_OPEN,
		}

		if result := <-api.Srv.Store.Team().Save(team, T); result.Err != nil {
			c.Err = result.Err
			return
		} else {

			createdTeam := result.Data.(*model.Team)

			channel := &model.Channel{DisplayName: "General", Name: "general", Type: model.CHANNEL_OPEN, TeamId: createdTeam.Id}
			if _, err := api.CreateChannel(c, channel, false, T); err != nil {
				c.Err = err
				return
			}

			teamID = createdTeam.Id
		}

		// Create user for testing
		user := &model.User{
			TeamId:   teamID,
			Email:    utils.RandomEmail(utils.Range{20, 20}, utils.LOWERCASE),
			Nickname: username[0],
			Password: api.USER_PASSWORD}

		result, err := client.CreateUser(user, "", T)
		if err != nil {
			c.Err = err
			return
		}
		api.Srv.Store.User().VerifyEmail(result.Data.(*model.User).Id, T)
		newuser := result.Data.(*model.User)
		userID = newuser.Id

		// Login as user to generate auth token
		_, err = client.LoginById(newuser.Id, api.USER_PASSWORD, T)
		if err != nil {
			c.Err = err
			return
		}

		// Respond with an auth token this can be overriden by a specific test as required
		sessionCookie := &http.Cookie{
			Name:     model.SESSION_COOKIE_TOKEN,
			Value:    client.AuthToken,
			Path:     "/",
			MaxAge:   model.SESSION_TIME_WEB_IN_SECS,
			HttpOnly: true,
		}
		http.SetCookie(w, sessionCookie)
		http.Redirect(w, r, "/channels/general", http.StatusTemporaryRedirect)
	}

	// Setup test environment
	env := TestEnvironment{
		Params:        params,
		Client:        client,
		CreatedTeamId: teamID,
		CreatedUserId: userID,
		Context:       c,
		Writer:        w,
		Request:       r,
	}

	// Grab the test ID and pick the test
	testname, ok := params["test"]
	var err2 *model.AppError
	switch testname[0] {
	case "autolink":
		err2 = testAutoLink(env)
		// ADD YOUR NEW TEST HERE!
	case "general":
		err2 = nil
	}

	if err != nil {
		c.Err = err2
		return
	}
}

func getChannelID(channelname string, teamid string, userid string) (id string, err bool) {
	// Grab all the channels
	T := i18n.GetSystemLanguage()
	result := <-api.Srv.Store.Channel().GetChannels(teamid, userid, T)
	if result.Err != nil {
		l4g.Debug("Unable to get channels")
		return "", false
	}

	data := result.Data.(*model.ChannelList)

	for _, channel := range data.Channels {
		if channel.Name == channelname {
			return channel.Id, true
		}
	}
	l4g.Debug("Could not find channel: " + channelname + ", " + strconv.Itoa(len(data.Channels)) + " possibilites searched")
	return "", false
}
