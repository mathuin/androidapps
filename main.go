package main

import (
	// "errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var host, port string

func main_env() {
	env_host, err := getenv("ANDROIDAPPS_HOST")
	if err != nil {
		log.Fatal(err)
	} else {
		host = env_host
	}

	env_port, err := getenv("ANDROIDAPPS_PORT")
	if err != nil {
		log.Fatal(err)
	} else {
		port = env_port
	}
}

func main_init() {
	// parse flags
}

func main() {
	// Check environment variables before flags
	main_env()
	database_env()
	web_env()

	// var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	// var ErrHelp = errors.New("flag: help requested")
	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// Now apply flags to settings
	main_init()
	database_init()
	web_init()

	var subcommand = flag.Arg(0)

	switch subcommand {
	case "runserver":
		database_init()
		web_init()
		main_init()
		http.HandleFunc("/", ServeIndex)
		http.HandleFunc("/static/", ServeStatic)
		http.HandleFunc("/media/", ServeMedia)
		log.Println("Starting server on host", host, "port", port)
		log.Fatal(http.ListenAndServe(host+":"+port, nil))
	default:
		Usage()
	}
}
