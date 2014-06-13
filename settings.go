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

func (s *Setting) set_value(key string) (err error) {
	if s.envvar != "" {
		s.value = os.Getenv(s.envvar)
	} else {
		err = fmt.Errorf("no environment variable for %s found", key)
	}
	if s.flag_value != "" {
		s.value = s.flag_value
	}
	if s.value == "" && err == nil {
		err = fmt.Errorf("no value for %s found -- set environment variable %s or use flag", key, s.envvar)
	}
	return
}

var settings map[string]*Setting

func init() {
	var err error
	settings = make(map[string]*Setting)
	settings["dbfile"], err = NewSetting("Database file", "ANDROIDAPPS_DBFILE")
	if err != nil {
		panic(err)
	}
	settings["host"], err = NewSetting("Hostname", "ANDROIDAPPS_HOST")
	if err != nil {
		panic(err)
	}
	settings["port"], err = NewSetting("Port", "ANDROIDAPPS_PORT")
	if err != nil {
		panic(err)
	}
	settings["name"], err = NewSetting("Developer name", "ANDROIDAPPS_NAME")
	if err != nil {
		panic(err)
	}
	settings["email"], err = NewSetting("Developer email address", "ANDROIDAPPS_EMAIL")
	if err != nil {
		panic(err)
	}

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
