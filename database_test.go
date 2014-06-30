package main

import (
	_ "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
)

var tempdir string

var database_tests = []struct {
	cmdargs []string
	output  string
}{
	{[]string{"reset"}, ""},
	{[]string{"list"}, "No apps are in the database!\n"},
	{[]string{"add", "./test/SimpleApp.apk"}, "The app simple.app was added!\n"},
	{[]string{"list"}, "  0: simple.app 1.0 SimpleApp 0\n"},
	{[]string{"enable", "simple.app"}, ""},
	{[]string{"list"}, "  0: simple.app 1.0 SimpleApp 1\n"},
	{[]string{"disable", "simple.app"}, ""},
	{[]string{"list"}, "  0: simple.app 1.0 SimpleApp 0\n"},
	{[]string{"remove", "simple.app"}, "The app simple.app was removed!\n"},
	{[]string{"list"}, "No apps are in the database!\n"},
}

func Test_database(t *testing.T) {
	// set environment variables
	var err error
	if tempdir == "" {
		tempdir, err = ioutil.TempDir("", "")
		checkErr(err, "TempDir failed")
	}
	os.Setenv("ANDROIDAPPS_DBFILE", path.Join(tempdir, "test.db"))
	os.Setenv("ANDROIDAPPS_HOST", "0.0.0.0")
	os.Setenv("ANDROIDAPPS_PORT", "4000")
	os.Setenv("ANDROIDAPPS_NAME", "Jane")
	os.Setenv("ANDROIDAPPS_EMAIL", "jane@example.com")
	// given this environment, let's execute a few commands!
	path, err := exec.LookPath("./androidapps")
	checkErr(err, "exec.LookPath failed")

	// run through all the commands
	for _, tt := range database_tests {
		ttcmd := exec.Command(path, tt.cmdargs...)
		ttout, err := ttcmd.Output()
		checkErr(err, "ttcmd.Output() failed")
		if string(ttout[:]) != tt.output {
			t.Errorf("Given cmdargs=%+v, wanted \"%+v\", got \"%+v\" instead", tt.cmdargs, tt.output, string(ttout[:]))
		}
	}
}
