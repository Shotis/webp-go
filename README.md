# webp-go
A Go package for working with libwebp

# Example

```go
func main() {
	// open our original PNG file
	file, err := os.Open("original.png")

	if err != nil {
		// end the program if it can't read the file
		log.Fatalln(err)
	}

	// decode the original image using PNG
	img, err := png.Decode(file)

	if err != nil {
		// end the program if it can't decode
		log.Fatalln(err)
	}

	// create a new webp.Config that can be converted 1:1 to a WebPConfig
	config := &Config{
		Lossless: true, // Use lossless compression
		Method:   0,    // use fastest compression. This ranges from 0-6
		Quality:  100,  // try preserving 100% quality
	}

	// create a new webp.Picture. webp.Picture is a wrapper for the libwebp WebPPicture
	picture := NewPicture(img)
	picture.Init()       // initialize and allocate the webp.Picture with the parameters from the passed in Image
	defer picture.Free() // free when the execution is complete

	var buf bytes.Buffer // allocate a new byte buffer where the data wil be written to

	picture.Encode(&buf, config) // Encode the image to the buffer

	ioutil.WriteFile("output.webp", buf.Bytes(), os.ModePerm) // write it to the file
}

```
