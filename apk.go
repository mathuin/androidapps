package main

// This file contains all the code required to manipulate the APKs.

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var namere, versionre, labelre, iconre, imgre *regexp.Regexp

func init() {
	namere = regexp.MustCompile("name='([^']*)'")
	versionre = regexp.MustCompile("versionName='([^']*)'")
	labelre = regexp.MustCompile("label='([^']*)'")
	iconre = regexp.MustCompile("icon='([^']*)'")
	imgre = regexp.MustCompile(".*\\.([a-z]*)$")
}

func extract_info(filename string) (name string, version string, label string, icon string) {
	// The correct way to extract this information requires writing a Go package which disassembles APKs.
	// That's hard.
	// For now, I'm just going to run "aapt dump badging <filename>" and extract what I need.

	path, err := exec.LookPath("aapt")
	checkErr(err, "exec.LookPath failed")
	aaptcmd := exec.Command(path, "dump", "badging", filename)
	out, err := aaptcmd.Output()
	checkErr(err, "aaptcmd.Output() failed")
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "package: ") {
			name = namere.FindStringSubmatch(line)[1]
			version = versionre.FindStringSubmatch(line)[1]
		} else if strings.HasPrefix(line, "application: ") {
			label = labelre.FindStringSubmatch(line)[1]
			icon = iconre.FindStringSubmatch(line)[1]
		}
	}
	return name, version, label, icon
}

func copy_files(filename string, label string, icon string) {
	// icon's name is "media/icons/<label>.<suffix>"
	imgsuffix := imgre.FindStringSubmatch(icon)[1]
	icondest := fmt.Sprintf("media/icons/%s.%s", label, imgsuffix)
	// apk's name is "media/products/<filename>"
	apkdest := fmt.Sprintf("media/products/%s", filepath.Base(filename))

	// copy icon from apk to icons directory
	r, err := zip.OpenReader(filename)
	checkErr(err, "zip.OpenReader() failed")

	for _, f := range r.File {
		if f.Name == icon {
			rc, err := f.Open()
			checkErr(err, "f.Open() failed")
			nf, err := os.Create(icondest)
			checkErr(err, "os.Create() failed")
			defer nf.Close()
			_, err = io.Copy(nf, rc)
			checkErr(err, "io.Copy() failed")
		}
	}
	r.Close()

	// copy apk to products directory
	err = cp(apkdest, filename)
	checkErr(err, "cp failed")
}

// https://gist.github.com/elazarl/5507969
func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()
	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}
	return d.Close()
}
