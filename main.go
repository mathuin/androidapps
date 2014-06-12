package main

import (
	// "errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var host, port, dbfile, name, email string

func main_init() {
	// parse flags
}

func apply_env_flag(envvar string, flagval string) (string, error) {
	var retval string
	var errval error
	envval, err := getenv(envvar)
	if err == nil {
		retval = envval
	}
	if flagval != "" {
		retval = flagval
	}
	if retval == "" {
		errval = fmt.Errorf("no value found -- set %s or use command line flags", envvar)
	}
	return retval, errval
}

func init() {
	init_funcs = append(init_funcs, main_init)
}

var init_funcs []func() // array of funcs

func main() {
	// Flags are only for things that are subcommand agnostic.
	// Arguments that remain apply to subcommands.

	// Define flags.
	var flag_host, flag_port, flag_dbfile, flag_name, flag_email string
	flag.StringVar(&flag_host, "host", "", "hostname")
	flag.StringVar(&flag_port, "port", "", "port")
	flag.StringVar(&flag_dbfile, "dbfile", "", "database filename")
	flag.StringVar(&flag_name, "name", "", "developer name")
	flag.StringVar(&flag_email, "email", "", "developer email address")

	// Parse flags!
	flag.Parse()

	// Apply results.
	var err error
	host, err = apply_env_flag("ANDROIDAPPS_HOST", flag_host)
	if err != nil {
		log.Fatal(err)
	}
	port, err = apply_env_flag("ANDROIDAPPS_PORT", flag_port)
	if err != nil {
		log.Fatal(err)
	}
	dbfile, err = apply_env_flag("ANDROIDAPPS_DBFILE", flag_dbfile)
	if err != nil {
		log.Fatal(err)
	}
	name, err = apply_env_flag("ANDROIDAPPS_NAME", flag_name)
	if err != nil {
		log.Fatal(err)
	}
	email, err = apply_env_flag("ANDROIDAPPS_EMAIL", flag_email)
	if err != nil {
		log.Fatal(err)
	}

	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	var subcommand = flag.Arg(0)
	var subcommand_func func()

	switch subcommand {
	case "runserver":
		subcommand_func = runserver
	case "export":
		// export the database to standard output
	case "import":
		// import database in export format
	case "check":
		// check for files corresponding to products
		// - optional arguments: verbose?  string match?
	case "rebuild":
		// for each product, re-extract title, version, icon.
		// - optional arguments: all? string match?
	case "list":
		// list products in database
		// - optional arguments: string match? enabled?
	case "enable":
		// enable product (will need flag added to database)
		// - optional arguments: all? one?
	case "disable":
		// disable product (will need flag added to database)
		// - optional arguments: all? one?
	case "add":
		// add product to database
		// stdin/out form?  WEB?!?
	case "remove":
		// remove product from database
		// - optional arguments: force? (i.e., skip are-you-sure)
	case "upgrade":
		// upgrade product in database (upload new APK)
		// default:
		// 	Usage()
	}

	if subcommand_func != nil {
		// Initialize submodules.
		for _, value := range init_funcs {
			value()
		}
		// Run the subcommand.
		subcommand_func()
	} else {
		Usage()
	}
}

func runserver() {
	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/static/", ServeStatic)
	http.HandleFunc("/media/", ServeMedia)
	log.Println("Starting server on host", host, "port", port)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}
