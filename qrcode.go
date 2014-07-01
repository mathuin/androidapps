package main

import (
	"fmt"
	"github.com/qpliu/qrencode-go/qrencode"
	"image/png"
	"os"
)

func make_qrcode(name, filename string) error {
	msg := fmt.Sprintf("market://search?q=pname:%s", name)
	grid, err := qrencode.Encode(msg, qrencode.ECLevelQ)
	checkErr(err, "Encode() failed")
	f, err := os.Create(filename)
	checkErr(err, "os.Create() failed")
	defer f.Close()
	png.Encode(f, grid.Image(8))
	return nil
}
