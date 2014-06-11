package main

import (
	"fmt"
	"os"
)

func getenv(envvar string) (value string, err error) {
	raw_value := os.Getenv(envvar)
	if raw_value == "" {
		return "", fmt.Errorf("%s has no value!", envvar)
	} else {
		return raw_value, nil
	}
}
