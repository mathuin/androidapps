package main

import (
	"flag"
	"fmt"
	"os"
)

// Setting is the struct representing a setting.
type Setting struct {
	description string // "hostname"
	envvar      string // environment variable
	flagValue   string // for flag.StringVar
	flagDefault string // for flag.StringVar
	value       string // actual value
}

// NewSetting creates a new setting.
func NewSetting(description, envvar string) (s *Setting, err error) {
	if description == "" {
		err = fmt.Errorf("missing description")
	} else if envvar == "" {
		err = fmt.Errorf("missing envvar")
	} else {
		s = &Setting{description: description, envvar: envvar}
	}
	return
}

func (s *Setting) setValue(key string) (err error) {
	if s.envvar != "" {
		s.value = os.Getenv(s.envvar)
	} else {
		err = fmt.Errorf("no environment variable for %s found", key)
	}
	if s.flagValue != "" {
		s.value = s.flagValue
	}
	if s.value == "" && err == nil {
		err = fmt.Errorf("no value for %s found -- set environment variable %s or use flag", key, s.envvar)
	}
	return
}

func applySettings(settings map[string]*Setting) {
	for key, s := range settings {
		err := s.setValue(key)
		checkErr(err, "set_value failed")
	}
}

var settings map[string]*Setting

func init() {
	var err error
	settings = make(map[string]*Setting)
	settings["dbfile"], err = NewSetting("Database file", "ANDROIDAPPS_DBFILE")
	checkErr(err, "setting dbfile failed")
	settings["host"], err = NewSetting("Hostname", "ANDROIDAPPS_HOST")
	checkErr(err, "setting host failed")
	settings["port"], err = NewSetting("Port", "ANDROIDAPPS_PORT")
	checkErr(err, "setting port failed")
	settings["name"], err = NewSetting("Developer name", "ANDROIDAPPS_NAME")
	checkErr(err, "setting name failed")
	settings["email"], err = NewSetting("Developer email address", "ANDROIDAPPS_EMAIL")
	checkErr(err, "setting email failed")

	// Define flags.
	for key, s := range settings {
		flag.StringVar(&s.flagValue, key, s.flagDefault, s.description)
	}
}
