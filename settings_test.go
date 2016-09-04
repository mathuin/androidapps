package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var NewSettingTests = []struct {
	desc   string
	envvar string
	result *Setting
	err    error
}{
	{"Port", "ANDROIDAPPS_PORT", &Setting{description: "Port", envvar: "ANDROIDAPPS_PORT"}, nil},
	{"Port", "", nil, fmt.Errorf("missing envvar")},
	{"", "ANDROIDAPPS_PORT", nil, fmt.Errorf("missing description")},
	{"", "", nil, fmt.Errorf("missing description")},
}

func Test_NewSetting(t *testing.T) {
	for _, tt := range NewSettingTests {
		output, err := NewSetting(tt.desc, tt.envvar)
		if !reflect.DeepEqual(output, tt.result) || !reflect.DeepEqual(err, tt.err) {
			t.Errorf("Given desc=%+#v and envvar=%+#v, wanted (%+#v, %+#v), got (%+#v, %+#v) instead", tt.desc, tt.envvar, tt.result, tt.err, output, err)
		}
	}
}

var SetValueTests = []struct {
	sin      *Setting // envvar and flag_value are the key values here
	envval   string
	checkval string
	err      error
}{
	{&Setting{}, "", "", fmt.Errorf("no environment variable for test found")},
	{&Setting{envvar: "TEST"}, "", "", fmt.Errorf("no value for test found -- set environment variable TEST or use flag")},
	{&Setting{envvar: "TEST"}, "foo", "foo", nil},
	{&Setting{envvar: "TEST", flagValue: "bar"}, "", "bar", nil},
	{&Setting{envvar: "TEST", flagValue: "bar"}, "foo", "bar", nil},
}

func Test_setValue(t *testing.T) {
	for _, tt := range SetValueTests {
		s := tt.sin
		if s.envvar != "" {
			oldenv := os.Getenv(s.envvar)
			os.Setenv(s.envvar, tt.envval)
			defer os.Setenv(s.envvar, oldenv)
		}
		err := s.setValue("test")
		if s.value != tt.checkval || !reflect.DeepEqual(err, tt.err) {
			t.Errorf("Given sin=%+#v and envval=%+#v, wanted (%+#v, %+#v), got (%+#v, %+#v) instead", tt.sin, tt.envval, tt.checkval, tt.err, s.value, err)
		}
	}
}

var applySettingsTests = []struct {
	settings map[string]*Setting
	envval   string
	checkval string
}{
	{map[string]*Setting{"test": &Setting{envvar: "TEST"}}, "new", "new"},
}

func Test_applySettings(t *testing.T) {
	for _, tt := range applySettingsTests {
		for _, s := range tt.settings {
			if s.envvar != "" {
				oldenv := os.Getenv(s.envvar)
				os.Setenv(s.envvar, tt.envval)
				defer os.Setenv(s.envvar, oldenv)
			}
		}
		applySettings(tt.settings)
		for _, s := range tt.settings {
			if s.value != tt.checkval {
				t.Errorf("Given s=%+#v and envval=%+#v, wanted %+#v, got %+#v", s, tt.envval, tt.checkval, s.value)
			}
		}
	}
}
