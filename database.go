package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type App struct {
	Id          int
	Title       string
	Apkfile     string
	Iconfile    string
	Description string
	Recent      string
	Package     string
	Version     string
}

var apps []App

func init() {
	dbfile, err := getenv("ANDROIDAPPS_DB")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select id, title, apkfile, iconfile, description, recent, package, version from store_product")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var a App
		rows.Scan(&a.Id, &a.Title, &a.Apkfile, &a.Iconfile, &a.Description, &a.Recent, &a.Package, &a.Version)
		apps = append(apps, a)
	}
}
