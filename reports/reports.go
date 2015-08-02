// Copyright 2015 Ondrej Donek. All rights reserved.
// See LICENSE file for more informations about licensing.

// Package for generating reports.
package reports

import (
	"database/sql"
//	"errors"
//	"fmt"
	"github.com/odTimeTracker/odtimetracker-go-lib/database"
	odstrings "github.com/odTimeTracker/odtimetracker-go-lib/strings"
	"log"
	"strconv"
	"time"
)

// Used for specifying report type.
type ReportType struct{ value string }

func (rt *ReportType) String() string { return rt.value }

// Types of report.
var ReportTypeDaily = &ReportType{value: "daily"}
var ReportTypeWeekly = &ReportType{value: "weekly"}
var ReportTypeMonthly = &ReportType{value: "monthly"}

// Used for specifying report output format.
type ReportFormatType struct{ value string }

func (rft *ReportFormatType) String() string { return rft.value }

// Available report formats.
var ReportFormatHtml = &ReportFormatType{value: "html"}

// Struct for reports.
type Report struct {
	Format    *ReportFormatType // Report format (HTML|PDF|...)
	ProjectId int64             // Filter report by project's ID (can be empty)
	Tags      string            // Filter report by tags (can be empty)
	Type      *ReportType       // Report type (daily|weekly|monthly)
	Output    string            // Output document
}

// Create new report.
func NewReport(db *sql.DB, rtype *ReportType, format *ReportFormatType, pid int64, tags string) (r Report) {
	r.Type = rtype
	r.Format = format
	r.ProjectId = pid
	r.Tags = tags

	// TODO Template file would be better solution...
	// TODO If we will use templates they can be user defined...
	r.Output = "<!DOCTYPE html>"
	r.Output += "<html lang=\"en\">"
	r.Output += "<head>"
	r.Output += "<meta charset=\"utf-8\">"
	// TODO Finish title according to used filter and report's type
	r.Output += "<title>odTimeTracker Report</title>"
	r.Output += "</head>"
	r.Output += "<body>"
	r.Output += "<div class=\"container\">"
	r.Output += "<h1>odTimeTracker Report</h1>"

	activities, err := database.SelectActivities(db, -1)
	if err != nil {
		log.Fatal(err)
	}

	// Summary per activities
	r.Output += "<h2>Activities</h2>"
	r.Output += "<p>Activities you have been working on.</p>"
	r.Output += "<table>"
	r.Output += "<thead>"
	r.Output += "<tr>"
	r.Output += "<th scope=\"col\">AID</th>"
	r.Output += "<th scope=\"col\">Project</th>"
	r.Output += "<th scope=\"col\">Name and description</th>"
	r.Output += "<th scope=\"col\">Tags</th>"
	r.Output += "<th scope=\"col\">Started</th>"
	r.Output += "<th scope=\"col\">Stopped</th>"
	r.Output += "<th scope=\"col\">Duration</th>"
	r.Output += "</tr>"
	r.Output += "</thead>"
	r.Output += "<tbody>"

	var totalDuration time.Duration

	for _, a := range activities {
		// TODO This needs to be rewritten!!!
		var p database.Project
		projects, _ := database.SelectProjectById(db, a.ProjectId)
		if len(projects) == 1 {
			p = projects[0]
		}

		starttime, err := a.StartedTime()
		if err != nil {
			starttime = time.Now()
		}
		started := odstrings.FormatTime(starttime)

		var stopped string
		stoptime, err := a.StoppedTime()
		if err != nil {
			stopped = ""
		} else {
			stopped = odstrings.FormatTime(stoptime)
		}

		var name string
		if a.Description != "" {
			name = a.Name + "<br><small>" + a.Description + "</small>"
		} else {
			name = a.Name
		}

		totalDuration += a.Duration()

		r.Output += "<tr>"
		r.Output += "<td>" + strconv.FormatInt(a.ActivityId, 10) + "</th>"
		r.Output += "<td>" + p.Name + "</th>"
		r.Output += "<td>" + name + "</th>"
		r.Output += "<td>" + a.Tags + "</th>"
		r.Output += "<td>" + started + "</th>"
		r.Output += "<td>" + stopped + "</th>"
		r.Output += "<td>" + a.Duration().String() + "</th>"
		r.Output += "</tr>"
	}

	r.Output += "</tbody>"
	r.Output += "<tfoot>"
	r.Output += "<td colspan=\"6\"><strong>Total duration:</strong></td>"
	r.Output += "<td><strong>" + totalDuration.String() + "</strong></td>"
	r.Output += "</tfoot>"
	r.Output += "</table>"

	// Summary per projects
	r.Output += "<h2>Projects</h2>"

	// Summary per tags
	r.Output += "<h2>Tags</h2>"

	r.Output += "</div>"
	r.Output += "</body>"
	r.Output += "</html>"

	return r
}

// Render report.
func (r *Report) Render() string {
	// ...
	return r.Output
}
