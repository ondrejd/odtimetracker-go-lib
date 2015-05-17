// Copyright 2015 Ondrej Donek. All rights reserved.
// See LICENSE file for more informations about licensing.

// Holds all what neccessary for odTimeTracker database - datatypes 
// `Activity` and `Project` and functions for manipulating data.
package database

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"strings"
	"time"
)

// Project name separator in activity string.
const projectNameSep = "@"

// Tags separator in activity string.
const tagsSep = ";"

// Description separator in activity string.
const descriptionSep = "#"

// Schema for our SQLite database
const schemaSql = `
CREATE TABLE Projects (
	ProjectId INTEGER PRIMARY KEY, 
	Name TEXT,
	Description TEXT,
	Created TEXT NOT NULL
);
CREATE TABLE Activities (
	ActivityId INTEGER PRIMARY KEY,
	ProjectId INTEGER NOT NULL,
	Name TEXT,
	Description TEXT,
	Tags TEXT,
	Started TEXT NOT NULL,
	Stopped TEXT NOT NULL DEFAULT '',
	FOREIGN KEY(ProjectId) REFERENCES Projects(ProjectId) 
);
PRAGMA user_version = 1;
`

// Project definition
type Project struct {
	ProjectId   int64  // Numeric identifier of the project.
	Name        string // Name of the project.
	Description string // Description of the project.
	Created     string // Datetime of creation of the project (formatted by RFC3339).
}

// Returns `Created` string as regular instance of `time.Time`.
func (p *Project) CreatedTime() (time.Time, error) {
	return time.Parse(p.Created, time.RFC3339)
}

// Activity definition
type Activity struct {
	ActivityId  int64
	ProjectId   int64
	Name        string
	Description string
	Tags        string
	Started     string
	Stopped     string
	project     Project
}

// Get project associated with the activity.
func (a *Activity) GetProject() Project {
	return a.project
}

// Set project associated with the activity.
func (a *Activity) SetProject(p Project) {
	a.ProjectId = p.ProjectId
	a.project = p
}

// Returns `Started` string as regular instance of `time.Time`.
func (a *Activity) StartedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, a.Started)
}

// Returns `Stopped` string as regular instance of `time.Time`.
func (a *Activity) StoppedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, a.Stopped)
}

// Returns duration of the activity.
func (a *Activity) Duration() time.Duration {
	stopped, err := a.StoppedTime()
	if err != nil {
		stopped = time.Now()
	}

	started, err := a.StartedTime()
	if err != nil {
		started = time.Now()
	}

	return stopped.Sub(started)
}

// Initialize activity from the given string.
func (a *Activity) Parse(db *sql.DB, activityString string) error {
	aStr := strings.Trim(activityString, " \n\t")
	if aStr == "" {
		return errors.New("Empty string given!")
	}

	if strings.Count(aStr, projectNameSep) > 1 ||
		strings.Count(aStr, tagsSep) > 1 ||
		strings.Count(aStr, descriptionSep) > 1 {

		return errors.New("Given activity string is not well formed!")
	}

	hasProjectName := strings.Contains(aStr, projectNameSep)
	hasTags := strings.Contains(aStr, tagsSep)
	hasDescription := strings.Contains(aStr, descriptionSep)
	projectName := ""

	if hasDescription == true {
		parts := strings.Split(aStr, descriptionSep)
		aStr = parts[0]
		a.Description = parts[1]
	}

	if hasTags == true {
		parts := strings.Split(aStr, tagsSep)
		aStr = parts[0]
		a.Tags = parts[1]
	}

	if hasProjectName == true {
		parts := strings.Split(aStr, projectNameSep)
		aStr = parts[0]
		projectName = parts[1]
	}

	a.Name = aStr

	if hasProjectName == true && projectName != "" {
		projects, err := SelectProjectByName(db, projectName)
		if err != nil {
			log.Fatal(err)
		}

		if len(projects) == 1 {
			a.ProjectId = projects[0].ProjectId
			a.SetProject(projects[0])
		} else if len(projects) == 0 {
			p, err := InsertProject(db, projectName, "")
			if err != nil {
				log.Fatal(err)
			}
			a.ProjectId = p.ProjectId
		}
	}

	return nil
}

// ====================================================================
// Public package functions

// Initialize storage.
func InitStorage(path string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	//defer db.Close() // Moved to `main.go`

	// Check if we are need to create schema
	ver, err := schemaVersion(db)
	if err != nil {
		return db, err
	}
	if ver < 1 {
		err := schemaCreate(db)
		if err != nil {
			return db, err
		}
	}

	return db, nil
}

// Insert new activity.
func InsertActivity(db *sql.DB, pid int64, name string, desc string, tags string) (a Activity, err error) {
	sqlStmt := `
	INSERT INTO Activities 
	(ProjectId, Name, Description, Tags, Started) 
	VALUES (?, ?, ?, ?, ?)
	`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	started := time.Now().Format(time.RFC3339)
	res, err := stmt.Exec(pid, name, desc, tags, started)
	if err != nil {
		log.Fatal(err)
	}

	aid, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	a.ActivityId = aid
	a.ProjectId = pid
	a.Name = name
	a.Description = desc
	a.Tags = tags
	a.Started = started

	return a, nil
}

// Insert new project.
func InsertProject(db *sql.DB, name string, desc string) (p Project, err error) {
	sqlStmt := `
	INSERT INTO Projects 
	(Name, Description, Created) 
	VALUES (?, ?, ?)
	`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	created := time.Now().Format(time.RFC3339)
	res, err := stmt.Exec(name, desc, created)
	if err != nil {
		log.Fatal(err)
	}

	pid, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	p.ProjectId = pid
	p.Name = name
	p.Description = desc
	p.Created = created

	return p, nil
}

// Remove activity(-ies) with given Id(s) form the database.
func RemoveActivity(db *sql.DB, id ...int64) (int, error) {
	// ...
	return 0, nil
}

// Remove project(s) with given Id(s) form the database.
func RemoveProject(db *sql.DB, id ...int64) (int, error) {
	// ...
	return 0, nil
}

// Return activities.
func SelectActivities(db *sql.DB, limit int) (activities []Activity, err error) {
	stmtSql := `SELECT * FROM Activities ORDER BY Started DESC LIMIT ?`
	stmt, err := db.Prepare(stmtSql)
	if err != nil {
		return activities, err
	}

	rows, err := stmt.Query(limit)
	if err != nil {
		return activities, err
	}

	defer rows.Close()
	for rows.Next() {
		var a Activity
		rows.Scan(&a.ActivityId, &a.ProjectId, &a.Name, &a.Description, &a.Tags, &a.Started, &a.Stopped)
		activities = append(activities, a)
	}
	rows.Close()

	return activities, nil
}

// Return activity(-ies) by given ID(s).
func SelectActivityById(db *sql.DB, id ...int64) (a []Activity, err error) {
	// ...
	return a, nil
}

// Return currently running activity.
func SelectActivityRunning(db *sql.DB) (a Activity, err error) {
	row := db.QueryRow(`SELECT * FROM Activities WHERE Stopped IS "" LIMIT 1`)
	err = row.Scan(&a.ActivityId, &a.ProjectId, &a.Name, &a.Description, &a.Tags, &a.Started, &a.Stopped)
	if err != nil {
		return a, err
	}

	return a, nil
}

// Return projects.
func SelectProjects(db *sql.DB, limit int) (p []Project, err error) {
	stmtSql := `SELECT * FROM Projects ORDER BY Name ASC LIMIT ?`
	stmt, err := db.Prepare(stmtSql)
	if err != nil {
		return p, err
	}

	rows, err := stmt.Query(limit)
	if err != nil {
		return p, err
	}

	defer rows.Close()
	return parseProjectsFromRows(rows)
}

// Return project(s) by given ID(s).
func SelectProjectById(db *sql.DB, id ...int64) (p []Project, err error) {
	var ids []string
	for _, id := range id {
		ids = append(ids, strconv.FormatInt(id, 10))
	}
	idsStr := strings.Join(ids, ", ")

	sqlStmt := "SELECT * FROM Projects WHERE Id IN (" + idsStr + ")"
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return p, err
	}

	defer rows.Close()
	return parseProjectsFromRows(rows)
}

// Return single project by given name(s).
func SelectProjectByName(db *sql.DB, name ...string) (projects []Project, err error) {
	namesStr := strings.Join(name, "\", \"")
	namesStr = "\"" + namesStr + "\""

	sqlStmt := "SELECT * FROM Projects WHERE Name IN (" + namesStr + ")"
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return projects, err
	}

	defer rows.Close()
	return parseProjectsFromRows(rows)
}

// Update activity in the database
// Return 1 if update was successfull otherwise 0.
func UpdateActivity(db *sql.DB, a Activity) (cnt int64, err error) {
	sqlStmt := `
	UPDATE Activities 
	SET
	ProjectId = ?, 
	Name = ?, 
	Description = ?, 
	Tags = ?, 
	Started = ?, 
	Stopped = ? 
	WHERE ActivityId = ? 
	`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(a.ProjectId, a.Name, a.Description, a.Tags,
		a.Started, a.Stopped, a.ActivityId)
	if err != nil {
		log.Fatal(err)
	}

	return res.RowsAffected()
}

// Update project in the database.
// Return 1 if update was successfull otherwise 0.
func UpdateProject(db *sql.DB, p Project) (cnt int64, err error) {
	sqlStmt := `
	UPDATE Projects 
	SET 
	Name = ?, 
	Description = ?, 
	Created = ? 
	WHERE ProjectId = ? 
	`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(p.Name, p.Description, p.Created, p.ProjectId)
	if err != nil {
		log.Fatal(err)
	}

	return res.RowsAffected()
}

// ====================================================================
// Some "internal" functions used in code above


// Returns schema version.
func schemaVersion(db *sql.DB) (int, error) {
	var user_version int = 0
	row := db.QueryRow("PRAGMA user_version;")
	err := row.Scan(&user_version)
	if err != nil {
		return 0, err
	}
	return user_version, nil
}

// Creates database schema.
func schemaCreate(db *sql.DB) error {
	_, err := db.Exec(schemaSql)
	if err != nil {
		return err
	}
	return nil
}

// Helper method that converting results rows
// into regular instances of Project object.
func parseProjectsFromRows(rows *sql.Rows) (projects []Project, err error) {
	for rows.Next() {
		var p Project
		rows.Scan(&p.ProjectId, &p.Name, &p.Description, &p.Created)
		projects = append(projects, p)
	}
	rows.Close()

	return projects, nil
}
