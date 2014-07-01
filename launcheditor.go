package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// not tested -- should be fun. :-(
func launcheditor(filename string) error {
	path, err := exec.LookPath("vi")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

var comment_header []byte
var line_terminator []byte
var space []byte

func init() {
	comment_header = []byte{'#'}
	line_terminator = []byte{'\n'}
	space = []byte{' '}
}

func createfile(text string) string {
	f, err := ioutil.TempFile("", "")
	checkErr(err, "ioutil.TempFile() failed")
	defer f.Close()
	if text != "" {
		for _, text := range split(text, 72) {
			line := [][]byte{comment_header, []byte(text), line_terminator}
			f.Write(bytes.Join(line, space))
		}
	}
	return f.Name()
}

func retrievestring(filename string) string {
	buf, err := ioutil.ReadFile(filename)
	checkErr(err, "ioutil.ReadFile failed")
	lines := []string{}
	for _, line := range bytes.Split(buf, line_terminator) {
		if !bytes.HasPrefix(line, comment_header) {
			lines = append(lines, string(line))
		}
	}
	return strings.Join(lines, " ")
}

// formatting stuff
// thanks to Egon Elgre and his post to golang-nuts !
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isbreak(b byte) bool {
	return b == ' ' || b == ',' || b == ':' || b == '.' || b == '-'
}

func findsplit(s string, near int) int {
	const RAGGEDNESS = 10
	bound := max(0, near-RAGGEDNESS)
	for i := near; i > bound; i -= 1 {
		if isbreak(s[i]) {
			return i + 1
		}
	}
	return min(near, len(s))
}

func split(s string, n int) []string {
	ss := make([]string, 0, len(s)/n+2)
	for len(s) > n {
		splitat := findsplit(s, n)
		ss = append(ss, s[:splitat])
		s = s[splitat:]
	}
	return append(ss, s)
}
