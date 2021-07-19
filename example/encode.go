package main

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/shotis/webp-go"
)

//Encoding example
func encode(input, output string, config *webp.Config) error {
	f, err := os.Open(input)

	if err != nil {
		return err
	}

	img, err := png.Decode(f)

	if err != nil {
		return err
	}

	var buf bytes.Buffer

	picture := webp.NewPicture(img)

	picture.Init()       // initialize the picture
	defer picture.Free() // free it

	if err = picture.EncodeTo(&buf, config); err != nil {
		return err
	}

	return ioutil.WriteFile(output, buf.Bytes(), os.ModePerm)
}
