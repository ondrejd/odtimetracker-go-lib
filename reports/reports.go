// Copyright 2015 Ondrej Donek. All rights reserved.
// See LICENSE file for more informations about licensing.

// Package for generating reports.
package reports

import (
//db "github.com/odTimeTracker/odtimetracker-go-lib/database"
//"errors"
)

/*const (
	ReportTypeDaily: "daily"
	ReportTypeWeekly: "weekly"
	ReportTypeMonthly: "monthly"
)*/

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
	Format    ReportFormatType
	ProjectId int64
	Tags      string
	Type      ReportType
}

// Create new report.
func NewReport(aType ReportType, aFormat ReportFormatType, aProjectId int64, aTags string) (r Report) {
	r.Type = aType
	r.Format = aFormat
	r.ProjectId = aProjectId
	r.Tags = aTags

	return r
}

// Render report.
func (r *Report) Render() string {
	// ...
	return ""
}
