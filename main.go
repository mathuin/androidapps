package main

import (
	"flag"
	"fmt"
	"os"
)

var init_funcs []func()

func main() {
	// Parse flags!
	flag.Parse()

	apply_settings()

	var subcommands map[string]func([]string)
	subcommands = map[string]func([]string){
		"runserver": runserver,
		"extract":   extract,
	}

	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	if subcommand := subcommands[flag.Arg(0)]; subcommand != nil {
		// Initialize submodules, then run subcommand.
		for _, value := range init_funcs {
			value()
		}
		subcommand(flag.Args())
	} else {
		Usage()
	}
}
