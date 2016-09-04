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

var extractInfoTests = []struct {
	name, version, label, icon, filename string
	err                                  error
}{
	{"org.twilley.android.firstapp", "1.0", "FirstApp", "res/drawable-mdpi-v4/ic_launcher.png", "./test/FirstApp10.apk", nil},
}

func Test_extractInfo(t *testing.T) {
	for _, tt := range extractInfoTests {
		name, version, label, icon, err := extractInfo(tt.filename)
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
