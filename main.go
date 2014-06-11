package main

import (
	"log"
	"net/http"
)

func main() {
	host, err := getenv("ANDROIDAPPS_HOST")
	if err != nil {
		log.Fatal(err)
	}
	port, err := getenv("ANDROIDAPPS_PORT")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/static/", ServeStatic)
	http.HandleFunc("/media/", ServeMedia)
	log.Println("Starting server on host", host, "port", port)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}
