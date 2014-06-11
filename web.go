package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func obfuscate(s string) template.HTML {
	var retval string
	for _, i := range s {
		retval += fmt.Sprintf("&#x%x;", i)
	}
	return template.HTML(retval)
}

func mailto(s string) template.HTMLAttr {
	return template.HTMLAttr(`href="` + obfuscate("mailto:"+s) + `"`)
}

type Page struct {
	Developer map[string]string
	Content   interface{}
}

var name, email string

func web_env() {
	env_name, err := getenv("ANDROIDAPPS_NAME")
	if err != nil {
		log.Fatal(err)
	} else {
		name = env_name
	}

	env_email, err := getenv("ANDROIDAPPS_EMAIL")
	if err != nil {
		log.Fatal(err)
	} else {
		email = env_email
	}
}

var appsPage *template.Template
var dev map[string]string

func web_init() {
	// parse flags

	dev = map[string]string{"name": name, "email": email}

	funcmap := template.FuncMap{
		"obfuscate": obfuscate,
		"mailto":    mailto,
	}
	layout := template.New("layout.html").Funcs(funcmap)
	layout = template.Must(layout.ParseFiles("templates/layout.html"))
	appsPage = template.Must(layout.Clone())
	appsPage = template.Must(appsPage.ParseFiles("templates/apps.html"))
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	err := appsPage.Execute(w, Page{
		Content:   apps,
		Developer: dev,
	})
	if err != nil {
		panic(err)
	}
}

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func ServeMedia(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}
