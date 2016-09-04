package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/qpliu/qrencode-go/qrencode"
)

func makeQRCode(name, filename string) error {
	msg := fmt.Sprintf("market://search?q=pname:%s", name)
	grid, err := qrencode.Encode(msg, qrencode.ECLevelQ)
	checkErr(err, "Encode() failed")
	f, err := os.Create(filename)
	checkErr(err, "os.Create() failed")
	defer f.Close()
	png.Encode(f, grid.Image(8))
	return nil
}
