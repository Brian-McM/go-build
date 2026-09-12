// Copyright (c) 2026 Tigera, Inc. All rights reserved.

package createvm

import (
	"strings"
	"testing"
)

func TestParseDiskGB(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int64
	}{
		{"200GB", 200},
		{"200G", 200},
		{"200", 200},
		{" 200gb ", 200},
		{"1500GB", 1500},
	} {
		got, err := parseDiskGB(tc.in)
		if err != nil {
			t.Errorf("parseDiskGB(%q) errored: %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("parseDiskGB(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestParseDiskGBRejectsGarbage(t *testing.T) {
	// A non-positive size parses fine but reaches the GCE API and fails there,
	// a long way from the env var that caused it.
	for _, in := range []string{"", "GB", "200GBx", "two hundred", "200TB", "0", "0GB", "-10GB", "-1"} {
		if _, err := parseDiskGB(in); err == nil {
			t.Errorf("parseDiskGB(%q) accepted garbage", in)
		}
	}
}

// The error must name the value the operator set, not the parser's leftovers --
// that is what makes a bad env var findable from a CI log.
func TestParseDiskGBErrorQuotesOriginalInput(t *testing.T) {
	_, err := parseDiskGB(" 200GBx ")
	if err == nil {
		t.Fatal("want an error")
	}
	if want := `" 200GBx "`; !strings.Contains(err.Error(), want) {
		t.Errorf("error %q does not contain the original input %s", err, want)
	}
}

// A non-positive duration parses fine but makes instanceSpec omit Scheduling
// entirely, leaving the VM with no reclaim deadline -- the backstop deletevm
// assumes when it returns 0 on a failed cleanup.
func TestParseMaxRun(t *testing.T) {
	for _, in := range []string{"90m", "1h", "30s", "1h30m"} {
		d, err := parseMaxRun(in)
		if err != nil || d <= 0 {
			t.Errorf("parseMaxRun(%q) = %v, %v", in, d, err)
		}
	}
	for _, in := range []string{"0", "0s", "-5m", "-1h", "", "ninety minutes", "90"} {
		if d, err := parseMaxRun(in); err == nil {
			t.Errorf("parseMaxRun(%q) accepted it, returning %v", in, d)
		}
	}
	if _, err := parseMaxRun("-5m"); err == nil || !strings.Contains(err.Error(), `"-5m"`) {
		t.Errorf("error should quote the input: %v", err)
	}
}
