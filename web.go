package main

import (
	"fmt"
	"html/template"
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

var layout *template.Template
var dev map[string]string

func web_init() {
	dev = map[string]string{"name": settings["name"].value, "email": settings["email"].value}

	funcmap := template.FuncMap{
		"obfuscate": obfuscate,
		"mailto":    mailto,
		"changes":   changes,
	}

	var err error
	layout, err = template.New("layout.html").Funcs(funcmap).ParseFiles("templates/layout.html", "templates/apps.html", "templates/changes.html")
	checkErr(err, "template.ParseFiles() failed")
}

func init() {
	init_funcs = append(init_funcs, web_init)
}

// not tested
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	// only enabled apps here
	apps := applist(true)
	err := layout.Execute(w, Page{
		Content:   apps,
		Developer: dev,
	})
	checkErr(err, "layout.Execute() failed")
}

// not tested
func ServeStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

// not tested
func ServeMedia(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

// not tested
func runserver(args []string) error {
	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/static/", ServeStatic)
	http.HandleFunc("/media/", ServeMedia)
	hostport := fmt.Sprintf("%s:%s", settings["host"].value, settings["port"].value)
	fmt.Println("Starting server on", hostport)
	return http.ListenAndServe(hostport, nil)
}
