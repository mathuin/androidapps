package main

import (
	_ "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	//"os/exec"
	"path"
	"testing"
)

var tempdir string

var database_tests = []struct {
	cmd    subcommand
	args   []string
	output string
}{
	{reset, []string{"reset"}, ""},
	{list, []string{"list"}, "No apps are in the database!\n"},
	{add, []string{"add", "./test/SimpleApp.apk"}, "The app simple.app was added!\n"},
	{list, []string{"list"}, "  0: simple.app 1.0 SimpleApp 0\n"},
	{enable, []string{"enable", "simple.app"}, ""},
	{list, []string{"list"}, "  0: simple.app 1.0 SimpleApp 1\n"},
	{disable, []string{"disable", "simple.app"}, ""},
	{list, []string{"list"}, "  0: simple.app 1.0 SimpleApp 0\n"},
	{remove, []string{"remove", "simple.app"}, "The app simple.app was removed!\n"},
	{list, []string{"list"}, "No apps are in the database!\n"},
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

	dbmap = initDb()
	defer dbmap.Db.Close()

	// run through all the commands
	for _, tt := range database_tests {
		tempFile, _ := ioutil.TempFile("", "stdout")
		oldStdout := os.Stdout
		os.Stdout = tempFile
		tt.cmd(tt.args)
		os.Stdout = oldStdout
		tempFile.Close()
		ttout, _ := ioutil.ReadFile(tempFile.Name())
		if string(ttout) != tt.output {
			t.Errorf("Given args=%+v, wanted \"%+v\", got \"%+v\" instead", tt.args, tt.output, string(ttout))
		}
	}
}
