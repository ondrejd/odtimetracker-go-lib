// Copyright 2015 Ondrej Donek. All rights reserved.
// See LICENSE file for more informations about licensing.

package strings

import (
	"testing"
	"time"
)

func TestReverse(t *testing.T) {
	cases := []struct {in, out string}{
		{"2015-05-08T05:05:33+02:00", "8.5 2015 05:05:33"},
		{"2015-05-08T20:22:30+02:00", "8.5 2015 20:22:30"},
		{"2015-05-16T23:13:31+02:00", "16.5 2015 23:13:31"},
	}
	for _, c := range cases {
		intime, _ := time.Parse(time.RFC3339, c.in)
		got := FormatTime(intime)
		if got != c.out {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.out)
		}
	}
}

func TestFormatTimeShort(t *testing.T) {
	cases := []struct {in, out string}{
		{"2015-05-08T05:05:33+02:00", "8.5 2015 5:05"},
		{"2015-05-08T20:22:30+02:00", "8.5 2015 20:22"},
		{"2015-05-16T23:13:31+02:00", "16.5 2015 23:13"},
	}
	for _, c := range cases {
		intime, _ := time.Parse(time.RFC3339, c.in)
		got := FormatTimeShort(intime)
		if got != c.out {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.out)
		}
	}
}
