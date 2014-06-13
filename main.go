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
		"list":      list,
		"enable":    enable,
		"disable":   disable,
		"add":       add,
		"remove":    remove,
		"reset":     reset,
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
		err := command(flag.Args())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		Usage()
	}
}
