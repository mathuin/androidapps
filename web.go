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

var appsPage *template.Template
var dev map[string]string

func web_init() {
	dev = map[string]string{"name": settings["name"].value, "email": settings["email"].value}

	funcmap := template.FuncMap{
		"obfuscate": obfuscate,
		"mailto":    mailto,
	}

	layout := template.New("layout.html").Funcs(funcmap)
	layout = template.Must(layout.ParseFiles("templates/layout.html"))
	appsPage = template.Must(layout.Clone())
	appsPage = template.Must(appsPage.ParseFiles("templates/apps.html"))
}

func init() {
	init_funcs = append(init_funcs, web_init)
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	apps := applist()
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

func runserver(args []string) error {
	var err error
	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/static/", ServeStatic)
	http.HandleFunc("/media/", ServeMedia)
	hostport := fmt.Sprintf("%s:%s", settings["host"].value, settings["port"].value)
	log.Println("Starting server on", hostport)
	err = http.ListenAndServe(hostport, nil)
	return err
}
