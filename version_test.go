// Copyright 2015 Ondrej Donek. All rights reserved.
// See LICENSE file for more informations about licensing.

package odtimetracker

import "testing"

const cv = "0.1.0"

func TestCurrentVersion(t *testing.T) {
	v := CurrentVersion()
	if cv != v.String() {
		t.Errorf("CurrentVersion == %s, expected %s", v.String(), cv)
	}
}

