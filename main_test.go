package main

import (
	"fmt"
	"reflect"
	"testing"
)

var exec_cmd_tests = []struct {
	// given args like "./androidapps reset", capture stdout, stderr, and err
	cmdargs []string
	stdout  string
	stderr  string
	err     error
}{
	{[]string{"reset"}, "", "", nil},
	{[]string{"boo"}, "", "", fmt.Errorf("bad args: ", []string{"boo"})},
}

func Test_exec_cmd(t *testing.T) {
	for _, tt := range exec_cmd_tests {
		err := exec_cmd(tt.cmdargs)
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("Given cmdargs=%+#v, wanted err %+#v, got err %+#v", tt.cmdargs, tt.err, err)
		}
	}
}
