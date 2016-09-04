package main

import (
	_ "database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var tempdir string

var databaseTests = []struct {
	cmd    subcommand
	args   []string
	stdout string
	stderr string
}{
	// regular sequence of events
	{reset, []string{"reset"}, "", ""},
	{list, []string{"list"}, "No apps are in the database!\n", ""},
	{upgrade, []string{"upgrade", "./test/FirstApp10.apk"}, "", "App org.twilley.android.firstapp does not already exist!"},
	{add, []string{"add", "./test/FirstApp10.apk", "-desc=pie"}, "The app org.twilley.android.firstapp was added!\n", ""},
	{add, []string{"add", "./test/FirstApp10.apk", "-desc=pie"}, "", "App org.twilley.android.firstapp already exists!"},
	// JMT: this needs to be fixed
	{upgrade, []string{"upgrade", "./test/FirstApp10.apk"}, "", "Cannot upgrade to existing version!"},
	{enable, []string{"enable", "org.twilley.android.firstapp"}, "The app org.twilley.android.firstapp was enabled!\n", ""},
	{enable, []string{"enable", "org.twilley.android.firstapp"}, "", "App org.twilley.android.firstapp was already enabled!"},
	{list, []string{"list"}, "Name:\n\torg.twilley.android.firstapp (enabled)\nVersion:\n\t1.0\nLabel:\n\tFirstApp\nDescription:\n\tpie\n", ""},
	{disable, []string{"disable", "org.twilley.android.firstapp"}, "The app org.twilley.android.firstapp was disabled!\n", ""},
	{disable, []string{"disable", "org.twilley.android.firstapp"}, "", "App org.twilley.android.firstapp was already disabled!"},
	{list, []string{"list"}, "Name:\n\torg.twilley.android.firstapp (not enabled)\nVersion:\n\t1.0\nLabel:\n\tFirstApp\nDescription:\n\tpie\n", ""},
	{remove, []string{"remove", "org.twilley.android.firstapp"}, "The app org.twilley.android.firstapp was removed!\n", ""},
	{list, []string{"list"}, "No apps are in the database!\n", ""},
	{remove, []string{"remove", "org.twilley.android.firstapp"}, "", "App org.twilley.android.firstapp does not exist!"},

	// bad arguments
	{reset, []string{"reset", "pie"}, "", "bad args: [reset pie]"},
	{add, []string{"add"}, "", "bad args: [add]"},
	{add, []string{"add", "./test/FirstApp10.apk", "pie"}, "", "bad args: [add ./test/FirstApp10.apk pie]"},
	{remove, []string{"remove"}, "", "bad args: [remove]"},
	{remove, []string{"remove", "org.twilley.android.firstapp", "pie"}, "", "bad args: [remove org.twilley.android.firstapp pie]"},
	{list, []string{"list", "pie"}, "", "bad args: [list pie]"},
	{enable, []string{"enable"}, "", "bad args: [enable]"},
	{enable, []string{"enable", "org.twilley.android.firstapp", "pie"}, "", "bad args: [enable org.twilley.android.firstapp pie]"},
	{disable, []string{"disable"}, "", "bad args: [disable]"},
	{disable, []string{"disable", "org.twilley.android.firstapp", "pie"}, "", "bad args: [disable org.twilley.android.firstapp pie]"},
	{upgrade, []string{"upgrade"}, "", "bad args: [upgrade]"},
	{upgrade, []string{"upgrade", "./test/FirstApp10.apk", "pie"}, "", "bad args: [upgrade ./test/FirstApp10.apk pie]"},
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
	for _, tt := range databaseTests {
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
