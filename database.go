package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"strings"
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
	Enabled     int64 // 0 = false, 1 = true
}

type Change struct {
	Id      int64
	Created int64
	Updated int64
	Name    string
	Ver     string
	Recent  string
}

// constructor
func newApp(name, ver, label, description string) App {
	return App{
		Created:     time.Now().UnixNano(),
		Name:        name,
		Ver:         ver,
		Label:       label,
		Description: description,
	}
}

func newChange(name, ver, recent string) Change {
	return Change{
		Created: time.Now().UnixNano(),
		Name:    name,
		Ver:     ver,
		Recent:  recent,
	}
}

func changes(name string) ([]Change, error) {
	if err := app_exists(name, func(a *App) error {
		return nil
	}); err == nil {
		// testing the simpler idea
		var changes []Change
		_, err := dbmap.Select(&changes, "select * from changes where name=?", name)
		checkErr(err, "Select failed")
		return changes, nil
	} else {
		return nil, fmt.Errorf("App %s does not exist!", name)
	}
}

func change(name, ver string) (myc Change, err error) {
	myc = Change{}
	err = dbmap.SelectOne(&myc, "select * from changes where name=? and ver=?", name, ver)
	return
}

func change_exists(name, ver string, cb func(c *Change) error) error {
	myc, err := change(name, ver)
	if err == nil {
		return cb(&myc)
	} else {
		return err
	}
}

func app(name string) (mya App, err error) {
	mya = App{}
	err = dbmap.SelectOne(&mya, "select * from apps where name=?", name)
	return
}

func app_exists(name string, cb func(a *App) error) error {
	mya, err := app(name)
	if err == nil {
		return cb(&mya)
	} else {
		return err
	}
}

// properly testing this requires good database fixtures
func applist(enabled bool) []App {
	var apps []App
	var selstr string
	if enabled == true {
		selstr = "select * from apps where enabled=1 order by id"
	} else {
		selstr = "select * from apps order by id"
	}
	_, err := dbmap.Select(&apps, selstr)
	checkErr(err, "Select failed")
	return apps
}

var dbmap *gorp.DbMap

func initDb() *gorp.DbMap {
	db, err := sql.Open("sqlite3", settings["dbfile"].value)
	checkErr(err, "sql.Open failed")

	mydbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// JMT: want names to be unique, no clue how
	mydbmap.AddTableWithName(App{}, "apps").SetKeys(true, "Id")
	mydbmap.AddTableWithName(Change{}, "changes").SetKeys(true, "Id").SetUniqueTogether("Name", "Ver")

	// JMT: eventually migrate/create elsewhere
	err = mydbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return mydbmap
}

// subcommands
// reset
func reset(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("bad args: %v", args)
	}
	return dbmap.TruncateTables()
}

const (
	add_header     string = `Please enter a description of the Android application.  Remember, this is what the customer will see when determining whether or not to install the software!`
	upgrade_header string = `Please describe the recent changes to the Android application.  Remember, this is what the customer will see when determining whether or not to install the software!`
)

// add
func add(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("bad args: %v", args)
	}
	filename := args[1]
	name, ver, label, icon, err := extract_info(filename)
	if err != nil {
		return err
	}
	if err := app_exists(name, func(a *App) error {
		return fmt.Errorf("App %s already exists!", name)
	}); err == sql.ErrNoRows {
		addflags := flag.NewFlagSet(args[0], flag.ExitOnError)
		descPtr := addflags.String("desc", "", "Description")

		addflags.Parse(args[2:])

		if len(addflags.Args()) > 0 {
			return fmt.Errorf("bad args: %v", args)
		}

		var desc string
		if *descPtr != "" {
			desc = *descPtr
		} else {
			// JMT: this code not tested!
			fpath := createfile(add_header, "")
			launcheditor(fpath)
			desc = retrievestring(fpath)
		}
		app := newApp(name, ver, label, desc)
		ierr := dbmap.Insert(&app)
		checkErr(ierr, "Insert failed")
		copy_files(filename, name, label, icon)
		fmt.Printf("The app %s was added!\n", name)
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
	if err := app_exists(name, func(a *App) error {
		_, derr := dbmap.Delete(a)
		return derr
	}); err == sql.ErrNoRows {
		return fmt.Errorf("App %s does not exist!", name)
	} else {
		fmt.Printf("The app %s was removed!\n", name)
		return err
	}
}

// list
func list(args []string) error {
	// JMT: consider using regular expressions or globs here
	if len(args) != 1 {
		return fmt.Errorf("bad args: %v", args)
	}
	// all apps
	apps := applist(false)
	if len(apps) == 0 {
		fmt.Println("No apps are in the database!")
	} else {
		for _, a := range apps {
			var enabled string
			if a.Enabled == 1 {
				enabled = "enabled"
			} else {
				enabled = "not enabled"
			}
			fmt.Printf("Name:\n\t%s (%s)\nVersion:\n\t%s\nLabel:\n\t%s\nDescription:\n", a.Name, enabled, a.Ver, a.Label)
			for _, line := range strings.Split(a.Description, string(line_terminator)) {
				fmt.Printf("\t%s\n", line)
			}
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
	return app_exists(name, func(a *App) error {
		if a.Enabled == 0 {
			if a.Description == "" {
				return fmt.Errorf("App %s has no description!", name)
			}
			a.Enabled = 1
			_, uerr := dbmap.Update(a)
			if uerr == nil {
				fmt.Printf("The app %s was enabled!\n", name)
			}
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
	return app_exists(name, func(a *App) error {
		if a.Enabled == 1 {
			a.Enabled = 0
			_, uerr := dbmap.Update(a)
			if uerr == nil {
				fmt.Printf("The app %s was disabled!\n", name)
			}
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
	name, ver, label, icon, err := extract_info(filename)
	if err != nil {
		return err
	}
	if err := app_exists(name, func(a *App) error {
		if ver == a.Ver {
			return fmt.Errorf("Cannot upgrade to existing version!")
		}
		if err := change_exists(name, ver, func(c *Change) error {
			return fmt.Errorf("Cannot upgrade to existing version!")
		}); err == nil {
			addflags := flag.NewFlagSet(args[0], flag.ExitOnError)
			recentPtr := addflags.String("recent", "", "Recent changes")

			addflags.Parse(args[2:])

			if len(addflags.Args()) > 0 {
				return fmt.Errorf("bad args: %v", args)
			}

			var recent string
			if *recentPtr != "" {
				recent = *recentPtr
			} else {
				fpath := createfile(upgrade_header, "")
				launcheditor(fpath)
				recent = retrievestring(fpath)
			}

			// begin transaction
			c := newChange(name, ver, recent)
			ierr := dbmap.Insert(c)
			checkErr(ierr, "Insert failed")

			copy_files(filename, name, label, icon)
			a.Updated = time.Now().UnixNano()
			a.Ver = ver
			a.Label = label
			_, uerr := dbmap.Update(a)
			if uerr == nil {
				fmt.Printf("The app %s was upgraded!\n", name)
			}
			// end transaction
			return uerr
		} else {
			return err
		}
	}); err == sql.ErrNoRows {
		return fmt.Errorf("App %s does not already exist!", name)
	} else {
		return err
	}
}

// modify
func modify(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("bad args: %v", args)
	}
	name := args[1]
	return app_exists(name, func(a *App) error {
		addflags := flag.NewFlagSet(args[0], flag.ExitOnError)
		descPtr := addflags.String("desc", "", "Description")
		recentPtr := addflags.String("recent", "", "Recent changes")

		addflags.Parse(args[2:])

		// JMT: if someone uses -recent on a description update (or
		// -desc on a recent changes update) they get the editor.  do
		// I care?

		remargs := addflags.Args()
		switch len(remargs) {
		case 0:
			// description
			var desc string
			if *descPtr != "" {
				desc = *descPtr
			} else {
				// JMT: this code not tested!
				fpath := createfile(add_header, a.Description)
				launcheditor(fpath)
				desc = retrievestring(fpath)
			}
			a.Description = desc
			_, uerr := dbmap.Update(a)
			if uerr == nil {
				fmt.Printf("The description for the app %s was modified successfully.\n", name)
			}
			return uerr
		case 1:
			// version
			ver := remargs[0]
			if err := change_exists(name, ver, func(c *Change) error {
				if *recentPtr != "" {
					c.Recent = *recentPtr
				} else {
					fpath := createfile(upgrade_header, c.Recent)
					launcheditor(fpath)
					c.Recent = retrievestring(fpath)
				}

				// begin transaction
				c.Updated = time.Now().UnixNano()
				_, uerr := dbmap.Update(c)
				checkErr(uerr, "Update failed")
				return nil
			}); err != nil {
				return err
			}
		default:
			// error
			return fmt.Errorf("bad args: %v", args)
		}
		return nil
	})
}
