package main

import (
	"fmt"
	"reflect"
	"testing"
)

var execCmdTests = []struct {
	// given args like "./androidapps reset", capture stdout, stderr, and err
	cmdargs []string
	stdout  string
	stderr  string
	err     error
}{
	{[]string{"reset"}, "", "", nil},
	{[]string{"boo"}, "", "", fmt.Errorf("bad args: %s", []string{"boo"})},
}

func Test_execCmd(t *testing.T) {
	for _, tt := range execCmdTests {
		err := execCmd(tt.cmdargs)
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("Given cmdargs=%+#v, wanted err %+#v, got err %+#v", tt.cmdargs, tt.err, err)
		}
	}
}
