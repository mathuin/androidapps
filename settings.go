package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Setting struct {
	description  string // "hostname"
	envvar       string // environment variable
	flag_value   string // for flag.StringVar
	flag_default string // for flag.StringVar
	value        string // actual value
}

func NewSetting(description, envvar string) *Setting {
	return &Setting{description: description, envvar: envvar}
}

func (s *Setting) set_value(key string) error {
	s.value = os.Getenv(s.envvar)
	if s.flag_value != "" {
		s.value = s.flag_value
	}
	if s.value == "" {
		return fmt.Errorf("no value for %s found -- set environment variable %s or use flag", key, s.envvar)
	}
	return nil
}

var settings map[string]*Setting

func init() {
	settings = make(map[string]*Setting)
	settings["dbfile"] = NewSetting("Database file", "ANDROIDAPPS_DBFILE")
	settings["host"] = NewSetting("Hostname", "ANDROIDAPPS_HOST")
	settings["port"] = NewSetting("Port", "ANDROIDAPPS_PORT")
	settings["name"] = NewSetting("Developer name", "ANDROIDAPPS_NAME")
	settings["email"] = NewSetting("Developer email address", "ANDROIDAPPS_EMAIL")

	// Define flags.
	for key, s := range settings {
		flag.StringVar(&s.flag_value, key, s.flag_default, s.description)
	}
}

func apply_settings() {
	for key, s := range settings {
		err := s.set_value(key)
		if err != nil {
			log.Fatal(err)
		}
	}
}
