package main

import (
	_ "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"path"
	"testing"
)

// for each fixture
// open tmpfile as database
// import appropriate fixture
// close database and delete file (with yields)

var tempdir string

func setup_database() {
	var err error
	if tempdir == "" {
		tempdir, err = ioutil.TempDir("", "")
		if err != nil {
			panic(err)
		}
	}
	settings = make(map[string]*Setting)
	settings["dbfile"], err = NewSetting("Database file", "TEST_DBFILE")
	settings["dbfile"].value = path.Join(tempdir, "testdb")
}

func Test_create_table(t *testing.T) {
	// These methods return no errors!

	setup_database()
	create_table()
	create_table()
}

func Test_drop_table(t *testing.T) {
	// These methods return no errors!

	setup_database()
	drop_table()
	drop_table()
}

// NB: the rest of these tests need a fixture.

func Test_is_enabled(t *testing.T) {
	setup_database()
	create_table()
}
