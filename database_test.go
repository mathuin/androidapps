package main

import (
	_ "database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

var tempdir string

var database_tests = []struct {
	cmd    subcommand
	args   []string
	stdout string
	stderr string
}{
	// regular sequence of events
	{reset, []string{"reset"}, "", ""},
	{list, []string{"list"}, "No apps are in the database!\n", ""},
	{upgrade, []string{"upgrade", "./test/SimpleApp.apk"}, "", "App simple.app does not already exist!"},
	{add, []string{"add", "./test/SimpleApp.apk", "-desc=pie"}, "The app simple.app was added!\n", ""},
	{add, []string{"add", "./test/SimpleApp.apk", "-desc=pie"}, "", "App simple.app already exists!"},
	// JMT: this needs to be fixed
	{upgrade, []string{"upgrade", "./test/SimpleApp.apk"}, "", "Cannot upgrade to existing version!"},
	{enable, []string{"enable", "simple.app"}, "The app simple.app was enabled!\n", ""},
	{enable, []string{"enable", "simple.app"}, "", "App simple.app was already enabled!"},
	{list, []string{"list"}, "Name:\n\tsimple.app (enabled)\nVersion:\n\t1.0\nLabel:\n\tSimpleApp\nDescription:\n\tpie\n", ""},
	{disable, []string{"disable", "simple.app"}, "The app simple.app was disabled!\n", ""},
	{disable, []string{"disable", "simple.app"}, "", "App simple.app was already disabled!"},
	{list, []string{"list"}, "Name:\n\tsimple.app (not enabled)\nVersion:\n\t1.0\nLabel:\n\tSimpleApp\nDescription:\n\tpie\n", ""},
	{remove, []string{"remove", "simple.app"}, "The app simple.app was removed!\n", ""},
	{list, []string{"list"}, "No apps are in the database!\n", ""},
	{remove, []string{"remove", "simple.app"}, "", "App simple.app does not exist!"},

	// bad arguments
	{reset, []string{"reset", "pie"}, "", "bad args: [reset pie]"},
	{add, []string{"add"}, "", "bad args: [add]"},
	{add, []string{"add", "./test/SimpleApp.apk", "pie"}, "", "bad args: [add ./test/SimpleApp.apk pie]"},
	{remove, []string{"remove"}, "", "bad args: [remove]"},
	{remove, []string{"remove", "simple.app", "pie"}, "", "bad args: [remove simple.app pie]"},
	{list, []string{"list", "pie"}, "", "bad args: [list pie]"},
	{enable, []string{"enable"}, "", "bad args: [enable]"},
	{enable, []string{"enable", "simple.app", "pie"}, "", "bad args: [enable simple.app pie]"},
	{disable, []string{"disable"}, "", "bad args: [disable]"},
	{disable, []string{"disable", "simple.app", "pie"}, "", "bad args: [disable simple.app pie]"},
	{upgrade, []string{"upgrade"}, "", "bad args: [upgrade]"},
	{upgrade, []string{"upgrade", "./test/SimpleApp.apk", "pie"}, "", "bad args: [upgrade ./test/SimpleApp.apk pie]"},
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
		var experr error
		if tt.stderr != "" {
			experr = fmt.Errorf(tt.stderr)
		}
		tempFileOut, _ := ioutil.TempFile("", "stdout")
		oldStdout := os.Stdout
		os.Stdout = tempFileOut
		acterr := tt.cmd(tt.args)
		os.Stdout = oldStdout
		tempFileOut.Close()
		ttout, _ := ioutil.ReadFile(tempFileOut.Name())
		if string(ttout) != tt.stdout || !reflect.DeepEqual(acterr, experr) {
			t.Errorf("Given args=%+v, wanted stdout \"%+v\" and stderr \"%+v\", got stdout \"%+v\" and stderr \"%+v\" instead", tt.args, tt.stdout, experr, string(ttout), acterr)
		}
	}
}
