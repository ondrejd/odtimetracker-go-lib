// Copyright 2015 Ondrej Donek. All rights reserved.
// See LICENSE file for more informations about licensing.

package odtimetracker

import "fmt"

// Struct defining simple version.
type Version struct {
	Major       int // Major version (the first part)
	Minor       int // Minor version (the second part)
	Maintenance int // Maintenance version (the third part)
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Maintenance)
}

// Current version of the library.
var currentVersion = &Version{
	Major: 0,
	Minor: 1,
	Maintenance: 0,
}

// Return current version of the library.
func CurrentVersion() *Version {
	return currentVersion
}

