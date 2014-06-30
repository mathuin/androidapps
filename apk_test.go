package main

import (
	"testing"
)

func check(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Wanted \"%+v\", got \"%+v\" instead", expected, actual)
	}
}

var extract_info_tests = []struct {
	name, version, label, icon, filename string
}{
	{"simple.app", "1.0", "SimpleApp", "res/drawable-mdpi/icon.png", "./test/SimpleApp.apk"},
}

func Test_extract_info(t *testing.T) {
	for _, tt := range extract_info_tests {
		name, version, label, icon := extract_info(tt.filename)
		check(t, tt.name, name)
		check(t, tt.version, version)
		check(t, tt.label, label)
		check(t, tt.icon, icon)
	}
}

// JMT: not sure how to test copy_files or cp right now...
