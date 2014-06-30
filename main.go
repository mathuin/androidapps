package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var init_funcs []func()

type subcommand func([]string) error

func main() {
	// Parse flags!
	flag.Parse()

	apply_settings()

	var subcommands map[string]subcommand
	subcommands = map[string]subcommand{
		"runserver": runserver,
		"reset":     reset,
		"add":       add,
		"remove":    remove,
		"list":      list,
		"enable":    enable,
		"disable":   disable,
		"upgrade":   upgrade,
	}

	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	if command := subcommands[flag.Arg(0)]; command != nil {
		// Initialize submodules, then run command.
		for _, value := range init_funcs {
			value()
		}

		// JMT: for now this is outside the hooks due to the defer
		dbmap = initDb()
		defer dbmap.Db.Close()

		err := command(flag.Args())
		checkErr(err, "command failed")

	} else {
		Usage()
	}
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
