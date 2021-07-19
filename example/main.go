package main

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/shotis/webp-go"
)

func main() {
	f, err := os.Open("example.png")

	if err != nil {
		panic(err)
	}

	img, err := png.Decode(f)

	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	picture := webp.NewPicture(img)
	config := &webp.Config{
		Lossless: true,
		Method:   0,
		Quality:  100,
	}

	picture.Init()       // initialize the picture
	defer picture.Free() // free it

	if err = picture.EncodeTo(&buf, config); err != nil {
		panic(err)
	}

	ioutil.WriteFile("example.webp", buf.Bytes(), os.ModePerm)
}
