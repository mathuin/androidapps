package main

// This file contains all the code required to manipulate the APKs.

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mathuin/axmlParser"
)

var imgre *regexp.Regexp

func init() {
	imgre = regexp.MustCompile(".*\\.([a-z]*)$")
}

func extractInfo(filename string) (name string, version string, label string, icon string, err error) {

	listener := new(axmlParser.AppNameListener)
	_, err = axmlParser.ParseApk(filename, listener)
	if err != nil {
		return
	}
	name = listener.PackageName
	version = listener.VersionName
	label = listener.ApplicationLabel
	icon = listener.ApplicationIcon

	return
}

func copyFiles(filename string, name string, label string, icon string, err error) {
	// icon's name is "media/icons/<label>.<suffix>"
	imgsuffix := imgre.FindStringSubmatch(icon)[1]
	icondest := fmt.Sprintf("media/icons/%s.%s", label, imgsuffix)
	// apk's name is "media/products/<filename>"
	apkdest := fmt.Sprintf("media/products/%s", filepath.Base(filename))

	// copy icon from apk to icons directory
	r, err := zip.OpenReader(filename)
	if err != nil {
		err = fmt.Errorf("zip.OpenReader() failed on archive %s", filename)
		return
	}

	for _, f := range r.File {
		if f.Name == icon {
			var rc io.ReadCloser
			rc, err = f.Open()
			if err != nil {
				err = fmt.Errorf("f.Open() failed on file %s", f.Name)
				return
			}
			var nf *os.File
			nf, err = os.Create(icondest)
			if err != nil {
				err = fmt.Errorf("os.Create() failed on file %s", icondest)
				return
			}
			defer nf.Close()
			_, err = io.Copy(nf, rc)
			checkErr(err, "io.Copy() failed")
		}
	}
	r.Close()

	// copy apk to products directory
	err = cp(apkdest, filename)
	checkErr(err, "cp failed")

	// generate QR code in target directory.
	qrdest := fmt.Sprintf("media/qrcodes/%s.%s", label, imgsuffix)
	err = makeQRCode(name, qrdest)
	checkErr(err, "make_qrcode() failed")
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
