// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"net/http"
	"strconv"

	l4g "github.com/alecthomas/log4go"
	"github.com/mattermost/mattermost-server/app"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/utils"
	"github.com/mssola/user_agent"
)

func (api *API) InitCompliance() {
	l4g.Debug(utils.T("api.compliance.init.debug"))

	api.BaseRoutes.Compliance.Handle("/reports", api.ApiSessionRequired(createComplianceReport)).Methods("POST")
	api.BaseRoutes.Compliance.Handle("/reports", api.ApiSessionRequired(getComplianceReports)).Methods("GET")
	api.BaseRoutes.Compliance.Handle("/reports/{report_id:[A-Za-z0-9]+}", api.ApiSessionRequired(getComplianceReport)).Methods("GET")
	api.BaseRoutes.Compliance.Handle("/reports/{report_id:[A-Za-z0-9]+}/download", api.ApiSessionRequiredTrustRequester(downloadComplianceReport)).Methods("GET")
}

func createComplianceReport(c *Context, w http.ResponseWriter, r *http.Request) {
	job := model.ComplianceFromJson(r.Body)
	if job == nil {
		c.SetInvalidParam("compliance")
		return
	}

	if !app.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	job.UserId = c.Session.UserId

	rjob, err := c.App.SaveComplianceReport(job)
	if err != nil {
		c.Err = err
		return
	}

	c.LogAudit("")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(rjob.ToJson()))
}

func getComplianceReports(c *Context, w http.ResponseWriter, r *http.Request) {
	if !app.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	crs, err := c.App.GetComplianceReports(c.Params.Page, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(crs.ToJson()))
}

func getComplianceReport(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireReportId()
	if c.Err != nil {
		return
	}

	if !app.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	job, err := c.App.GetComplianceReport(c.Params.ReportId)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(job.ToJson()))
}

func downloadComplianceReport(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireReportId()
	if c.Err != nil {
		return
	}

	if !app.SessionHasPermissionTo(c.Session, model.PERMISSION_MANAGE_SYSTEM) {
		c.SetPermissionError(model.PERMISSION_MANAGE_SYSTEM)
		return
	}

	job, err := c.App.GetComplianceReport(c.Params.ReportId)
	if err != nil {
		c.Err = err
		return
	}

	reportBytes, err := app.GetComplianceFile(job)
	if err != nil {
		c.Err = err
		return
	}

	c.LogAudit("downloaded " + job.Desc)

	w.Header().Set("Cache-Control", "max-age=2592000, public")
	w.Header().Set("Content-Length", strconv.Itoa(len(reportBytes)))
	w.Header().Del("Content-Type") // Content-Type will be set automatically by the http writer

	// attach extra headers to trigger a download on IE, Edge, and Safari
	ua := user_agent.New(r.UserAgent())
	bname, _ := ua.Browser()

	w.Header().Set("Content-Disposition", "attachment;filename=\""+job.JobName()+".zip\"")

	if bname == "Edge" || bname == "Internet Explorer" || bname == "Safari" {
		// trim off anything before the final / so we just get the file's name
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Write(reportBytes)
}
