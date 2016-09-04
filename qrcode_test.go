package main

import (
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

// given a package name "simple.app",
// build a qr code with string 'market://search?q=pname:simple.app'
// save as png with name "whatever.png"

func Test_makeQRCode(t *testing.T) {
	var err error
	if tempdir == "" {
		tempdir, err = ioutil.TempDir("", "")
		checkErr(err, "TempDir failed")
	}
	qrdest := path.Join(tempdir, "test.png")
	err = makeQRCode("org.twilley.android.firstapp", qrdest)
	if err != nil {
		t.Errorf("Dammit!")
	}

	expf, err := os.Open("./test/FirstApp-qrcode.png")
	checkErr(err, "os.Open() failed")
	expimg, _, err := image.Decode(expf)
	checkErr(err, "image.Decode(expf) failed")

	actf, err := os.Open(qrdest)
	checkErr(err, "os.Open() failed")
	actimg, _, err := image.Decode(actf)
	checkErr(err, "image.Decode(actf) failed")

	if !reflect.DeepEqual(expimg, actimg) {
		t.Errorf("expimage and actimg do not match")
	}

}
