// Copyright 2015 Ondrej Donek. All rights reserved.
// See LICENSE file for more informations about licensing.

package database

import (
	"testing"
)

func TestActivityParse(t *testing.T) {
	// TODO Init storage in memory!
	db, err := InitStorage(":memory:")
	if err != nil {
		t.Error("Expected nil got ", err)
	}
	defer db.Close()

	aStr1 := "Test activity;tag1,tag2"
	var a1 Activity
	err = a1.Parse(db, aStr1)
	if err != nil {
		t.Error("Expected nil got ", err)
	}

	if a1.Name != "Test activity" {
		t.Error("Expected 'Test activity' got ", a1.Name)
	}

	if a1.Tags != "tag1,tag2" {
		t.Error("Expected 'tag1,tag2' got ", a1.Tags)
	}
}
