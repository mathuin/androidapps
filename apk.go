package main

// This file contains all the code required to manipulate the APKs.

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
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
	if err != nil {
		log.Fatal(err)
	}
	aaptcmd := exec.Command(path, "dump", "badging", filename)
	out, err := aaptcmd.Output()
	if err != nil {
		log.Fatal(err)
	}
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
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range r.File {
		if f.Name == icon {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			f, err := os.Create(icondest)
			if err != nil {
				log.Fatal(err)
			}
			if _, err := io.Copy(f, rc); err != nil {
				f.Close()
			}
		}
	}
	r.Close()

	// copy apk to products directory
	err = cp(apkdest, filename)
	if err != nil {
		log.Fatal(err)
	}
}

// https://gist.github.com/elazarl/5507969
func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
