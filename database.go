package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type App struct {
	Name        string
	Version     string
	Label       string
	Description string
	Recent      string
	Enabled     bool
}

var apps map[string]App

func do(cb func(db *sql.DB)) {
	mydb, err := sql.Open("sqlite3", settings["dbfile"].value)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mydb.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	cb(mydb)
}

// create and drop table
func create_table() {
	do(func(db *sql.DB) {
		if _, err := db.Exec("create table if not exists apps (name varchar(100) primary key, version varchar(20), label varchar(20), description text, recent text, enabled bool)"); err != nil {
			log.Fatal(err)
		}
	})
}

func drop_table() {
	do(func(db *sql.DB) {
		if _, err := db.Exec("drop table if exists apps"); err != nil {
			log.Fatal(err)
		}
	})
}

// refresh global variable
// all apps
func refresh_apps() {
	do(func(db *sql.DB) {
		rows, err := db.Query("select name, version, label, description, recent, enabled from apps")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		apps = nil
		for rows.Next() {
			var a App
			rows.Scan(&a.Name, &a.Version, &a.Label, &a.Description, &a.Recent, &a.Enabled)
			// JMT: WTF
			if apps == nil {
				apps = make(map[string]App)
			}
			apps[a.Name] = a
		}
	})
}

func app_exists(name string) bool {
	var count int
	do(func(db *sql.DB) {
		check := db.QueryRow("select count(*) from apps where name = ?", name)
		if err := check.Scan(&count); err != nil {
			log.Fatal(err)
		}
	})
	if count == 1 {
		return true
	} else {
		return false
	}
}

func is_enabled(name string) (bool, error) {
	var enabled bool
	var err error
	switch app_exists(name) {
	case true:
		do(func(db *sql.DB) {
			check := db.QueryRow("select enabled from apps where name = ?", name)
			if err := check.Scan(&enabled); err != nil {
				log.Fatal(err)
			}
		})
	case false:
		err = fmt.Errorf("app %s does not exist", name)
	}
	return enabled, err
}

func refresh_app(name string) {
	switch app_exists(name) {
	case true:
		do(func(db *sql.DB) {
			row := db.QueryRow("select name, version, label, description, recent, enabled from apps where name = ?", name)
			var a App
			if err := row.Scan(&a.Name, &a.Version, &a.Label, &a.Description, &a.Recent, &a.Enabled); err != nil {
				log.Fatal(err)
			}
			// JMT: WTF
			if apps == nil {
				apps = make(map[string]App)
			}
			apps[name] = a
		})
	case false:
		delete(apps, name)
	}
}

// add record
func add_app(name string, version string, label string, description string, recent string, enabled int) error {
	var err error
	switch app_exists(name) {
	case false:
		do(func(db *sql.DB) {
			_, err := db.Exec("insert into apps (name, version, label, description, recent, enabled) values (?, ?, ?, ?, ?, ?)", name, version, label, description, recent, enabled)
			if err != nil {
				log.Fatal(err)
			}
		})
		refresh_app(name)
	case true:
		err = fmt.Errorf("app %s already exists", name)
	}
	return err
}

// delete record
func del_app(name string) {
	do(func(db *sql.DB) {
		_, err := db.Exec("delete from apps where name = ?", name)
		if err != nil {
			log.Fatal(err)
		}
	})
	refresh_app(name)
}

// modify record -- too hard at the moment

// enable record
func enable_app(name string) {
	do(func(db *sql.DB) {
		_, err := db.Exec("update apps set enabled = 1 where name = ?", name)
		if err != nil {
			log.Fatal(err)
		}
	})
	refresh_app(name)
}

// disable record
func disable_app(name string) {
	do(func(db *sql.DB) {
		_, err := db.Exec("update apps set enabled = 0 where name = ?", name)
		if err != nil {
			log.Fatal(err)
		}
	})
	refresh_app(name)
}

func database_init() {
	create_table()
	refresh_apps()
}

func init() {
	if apps == nil {
		apps = make(map[string]App)
	}
	init_funcs = append(init_funcs, database_init)
}

// subcommands
func reset(args []string) error {
	var err error
	for key := range apps {
		delete(apps, key)
	}
	drop_table()
	log.Println("The database was successfully reset!")
	return err
}

func add(args []string) error {
	var err error
	if len(args) != 2 {
		return fmt.Errorf("bad args: %v", args)
	}
	filename := args[1]
	name, version, label, icon := extract_info(filename)
	switch app_exists(name) {
	case false:
		copy_files(filename, label, icon)
		// JMT: need elegant solution for description
		add_app(name, version, label, "Description", "", 0)
		log.Printf("The app %s was successfully added from %s\n", name, filename)
	case true:
		err = fmt.Errorf("cannot add %s: %s already exists!", filename, name)
	}
	return err
}

func remove(args []string) error {
	var err error
	if len(args) != 2 {
		return fmt.Errorf("bad args: %v", args)
	}
	name := args[1]
	switch app_exists(name) {
	case true:
		// JMT: currently not deleting files
		del_app(name)
		log.Printf("The app %s was successfully removed!\n", name)
	case false:
		err = fmt.Errorf("App %s does not exist!", name)
	}
	return err
}

func list(args []string) error {
	var err error
	log.Println("List of apps:")
	for key := range apps {
		log.Printf("%+v\n", apps[key])
	}
	return err
}

func enable(args []string) error {
	var err error
	if len(args) != 2 {
		err = fmt.Errorf("bad args: %v", args)
		return err
	}
	name := args[1]
	check, err := is_enabled(name)
	if err != nil {
		log.Fatal(err)
	}
	switch check {
	case false:
		enable_app(name)
		log.Printf("The app %s was successfully enabled!\n", name)
	case true:
		err = fmt.Errorf("App %s was already enabled!", name)
	}
	return err
}

func disable(args []string) error {
	var err error
	if len(args) != 2 {
		return fmt.Errorf("bad args: %v", args)
	}
	name := args[1]
	check, err := is_enabled(name)
	if err != nil {
		log.Fatal(err)
	}
	switch check {
	case true:
		disable_app(name)
		log.Printf("The app %s was successfully disabled!\n", name)
		return nil
	case false:
		return fmt.Errorf("App %s was already disabled!", name)
	}
	log.Printf("The app %s was successfully disabled!\n", name)
	return err
}

func upgrade(args []string) error {
	var err error
	if len(args) != 2 {
		return fmt.Errorf("bad args: %v", args)
	}
	filename := args[1]
	name, version, label, icon := extract_info(filename)
	switch app_exists(name) {
	case true:
		copy_files(filename, label, icon)
		// JMT: need elegant solution for "recent" here
		olddesc = apps[name].Description
		del_app(name)
		add_app(name, version, label, olddesc, "Recent", 0)
		log.Printf("The app %s was successfully upgraded from %s\n", name, filename)
	case false:
		err = fmt.Errorf("App %s does not exist, use 'add' instead!", name)
	}
	return err
}
