// Copyright 2015 Ondrej Donek. All rights reserved.
// See LICENSE file for more informations about licensing.

// Some strings related utilities for odTimeTracker
package strings

import (
	"fmt"
	"time"
)

// Return time formated like this "12.1 2015 03:45:00".
func FormatTime(t time.Time) string {
	return fmt.Sprintf("%d.%d %d %02d:%02d:%02d", t.Day(), t.Month(), t.Year(),
		t.Hour(), t.Minute(), t.Second())
}

// Return time formated like this "12.1 2015 3:45".
func FormatTimeShort(t time.Time) string {
	return fmt.Sprintf("%d.%d %d %d:%02d", t.Day(), t.Month(), t.Year(),
		t.Hour(), t.Minute())
}
