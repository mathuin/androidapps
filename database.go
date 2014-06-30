package main

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

// JMT: should be unique on name, not sure how to do that.
type App struct {
	Id          int64
	Created     int64
	Updated     int64
	Name        string
	Ver         string
	Label       string
	Description string
	Recent      string
	Enabled     int64 // 0 = false, 1 = true
}

// constructor
func newApp(name, ver, label, description, recent string, enabled int64) App {
	return App{
		Created:     time.Now().UnixNano(),
		Name:        name,
		Ver:         ver,
		Label:       label,
		Description: description,
		Recent:      recent,
		Enabled:     enabled,
	}
}

func exists(name string, cb func(a *App) error) error {
	mya := App{}
	err := dbmap.SelectOne(&mya, "select * from apps where name=?", name)
	if err == nil {
		log.Println("app exists!")
		return cb(&mya)
	} else {
		return err
	}
}

func applist() []App {
	var apps []App
	_, err := dbmap.Select(&apps, "select * from apps order by id")
	checkErr(err, "Select failed")
	return apps
}

var dbmap *gorp.DbMap

func initDb() *gorp.DbMap {
	db, err := sql.Open("sqlite3", settings["dbfile"].value)
	checkErr(err, "sql.Open failed")

	mydbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	mydbmap.AddTableWithName(App{}, "apps").SetKeys(true, "Id")

	// JMT: eventually migrate/create elsewhere
	err = mydbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return mydbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

// subcommands
// reset
func reset(args []string) error {
	return dbmap.TruncateTables()
}

// add
func add(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("bad args: %v", args)
	}
	filename := args[1]
	name, ver, label, icon := extract_info(filename)
	if err := exists(name, func(a *App) error {
		return fmt.Errorf("App %s already exists!", name)
	}); err == sql.ErrNoRows {
		copy_files(filename, label, icon)
		// JMT: Description here!
		log.Printf("in add: got this far!")
		app := newApp(name, ver, label, "Description", "", int64(0))
		log.Printf("in add: %+v", app)
		ierr := dbmap.Insert(&app)
		checkErr(ierr, "Insert failed")
		return ierr
	} else {
		return err
	}
}

// remove
func remove(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("bad args: %v", args)
	}
	name := args[1]
	if err := exists(name, func(a *App) error {
		_, derr := dbmap.Delete(a)
		return derr
	}); err == sql.ErrNoRows {
		return fmt.Errorf("App %s does not exist!", name)
	} else {
		return err
	}
}

// list
func list(args []string) error {
	apps := applist()
	if len(apps) == 0 {
		log.Println("No apps are in the database!")
	} else {
		for x, a := range apps {
			log.Printf("  %d: %v\n", x, a)
		}
	}
	return nil
}

// enable
func enable(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("bad args: %v", args)
	}
	name := args[1]
	return exists(name, func(a *App) error {
		if a.Enabled == 0 {
			a.Enabled = 1
			_, uerr := dbmap.Update(a)
			return uerr
		} else {
			return fmt.Errorf("App %s was already enabled!", name)
		}
	})
}

// disable
func disable(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("bad args: %v", args)
	}
	name := args[1]
	return exists(name, func(a *App) error {
		if a.Enabled == 1 {
			a.Enabled = 0
			_, uerr := dbmap.Update(a)
			return uerr
		} else {
			return fmt.Errorf("App %s was already disabled!", name)
		}
	})
}

// upgrade
func upgrade(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("bad args: %v", args)
	}
	filename := args[1]
	name, ver, label, icon := extract_info(filename)
	if err := exists(name, func(a *App) error {
		copy_files(filename, label, icon)
		a.Updated = time.Now().UnixNano()
		a.Ver = ver
		a.Label = label
		// JMT: editor here!
		a.Recent = "Recent"
		_, uerr := dbmap.Update(a)
		return uerr
	}); err == sql.ErrNoRows {
		return fmt.Errorf("App %s does not already exist!", name)
	} else {
		return err
	}
}
