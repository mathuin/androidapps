package main

import (
	"reflect"
	"testing"
)

func check(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Wanted \"%+v\", got \"%+v\" instead", expected, actual)
	}
}

var extract_info_tests = []struct {
	name, version, label, icon, filename string
	err                                  error
}{
	{"simple.app", "1.0", "SimpleApp", "res/drawable-mdpi/icon.png", "./test/SimpleApp.apk", nil},
}

func Test_extract_info(t *testing.T) {
	for _, tt := range extract_info_tests {
		name, version, label, icon, err := extract_info(tt.filename)
		check(t, tt.name, name)
		check(t, tt.version, version)
		check(t, tt.label, label)
		check(t, tt.icon, icon)
		if !reflect.DeepEqual(tt.err, err) {
			t.Errorf("Wanted \"%+v\", got \"%+v\" instead", tt.err, err)
		}
	}
}

// JMT: not sure how to test copy_files or cp right now...
