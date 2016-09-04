package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var initFuncs []func()

type subcommand func([]string) error

var subcommands map[string]subcommand

// Usage has not yet been tested
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	subcommands = map[string]subcommand{
		"runserver": runserver,
		"reset":     reset,
		"add":       add,
		"remove":    remove,
		"list":      list,
		"enable":    enable,
		"disable":   disable,
		"upgrade":   upgrade,
		"modify":    modify,
	}
}

func execCmd(args []string) error {
	if command := subcommands[args[0]]; command != nil {
		// Initialize submodules, then run command.
		for _, value := range initFuncs {
			value()
		}

		// JMT: for now this is outside the hooks due to the defer
		dbmap = initDb()
		defer dbmap.Db.Close()

		err := command(args)
		checkErr(err, "command failed")
		return nil
	}
	return fmt.Errorf("bad args: %s", args)
}

// not tested
func main() {
	// Parse flags!
	flag.Parse()

	applySettings(settings)

	if err := execCmd(flag.Args()); err != nil {
		Usage()
	}
}

// not tested
func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
