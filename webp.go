package webp

/*
#cgo pkg-config: libwebp
#include <stdlib.h>
#include <webp/encode.h>

// A custom writer function for intercepting WebP data
int writeTo(uint8_t* data, size_t length, WebPPicture* picture);
*/
import "C"
import (
	"fmt"
	"image"
	"io"
	"unsafe"
)

type Config struct {
	Lossless bool
	Method   int
	Quality  float32

	cConf *C.struct_WebPConfig
}

func (c *Config) toCStruct() (*C.struct_WebPConfig, error) {
	if c.cConf != nil {
		return c.cConf, nil
	}

	var lossless int
	if c.Lossless {
		lossless = 1
	}

	var config C.struct_WebPConfig

	i, err := C.WebPConfigPreset(&config, C.WEBP_PRESET_PHOTO, C.float(c.Quality))

	if err != nil {
		return nil, err
	}

	if i != 1 {
		return nil, fmt.Errorf("error occurred initializing config: %d", i)
	}

	config.lossless = C.int(lossless)
	config.quality = C.float(c.Quality)
	config.method = C.int(c.Method)

	c.cConf = &config
	return &config, nil
}

type Picture struct {
	Width, Height int
	Stride        int
	Pixels        []uint8

	cpicture *C.struct_WebPPicture
}

func NewPicture(img image.Image) *Picture {
	switch p := img.(type) {
	case *image.RGBA:
		return &Picture{
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
			Pixels: p.Pix,
			Stride: p.Stride,
		}
	case *image.NRGBA:
		return &Picture{
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
			Pixels: p.Pix,
			Stride: p.Stride,
		}
	}
	return nil
}

func (p *Picture) Init() error {
	var pic C.struct_WebPPicture
	i, err := C.WebPPictureInit(&pic)
	if err != nil {
		return err
	}
	if i != 1 {
		return fmt.Errorf("error occurred initializing the picutre: %d", i)
	}

	pic.width = C.int(p.Width)
	pic.height = C.int(p.Height)

	err = p.alloc()

	if err != nil {
		return err
	}

	p.cpicture = &pic

	return nil
}

func (p *Picture) alloc() error {
	i, err := C.WebPPictureAlloc(p.cpicture)

	if err != nil {
		return err
	}

	if i != 1 {
		return fmt.Errorf("error occurred allocating picture: %d", i)
	}
	return nil
}

func (p *Picture) Free() error {
	_, err := C.WebPPictureFree(p.cpicture)

	if err != nil {
		return err
	}
	return nil
}

var writerMap map[unsafe.Pointer]io.Writer = make(map[unsafe.Pointer]io.Writer)

//export writeTo
func writeTo(data *C.uint8_t, data_size C.size_t, picture *C.struct_WebPPicture) C.int {
	w, ok := writerMap[unsafe.Pointer(picture)]

	if ok {
		read, err := w.Write(C.GoBytes(unsafe.Pointer(data), C.int(data_size)))

		if err != nil {
			return 0
		}

		return C.int(read)
	}

	return C.int(0)
}

func (p *Picture) Encode(w io.Writer, config *Config) error {
	conf, err := config.toCStruct()

	if err != nil {
		return err
	}

	p.cpicture.writer = C.WebPWriterFunction(C.writeTo)
	writerMap[unsafe.Pointer(p.cpicture)] = w

	defer func() {
		delete(writerMap, unsafe.Pointer(p.cpicture))
	}()

	cb := C.CBytes([]byte(p.Pixels))
	defer C.free(cb)

	// import our pixels into the image
	status, err := C.WebPPictureImportRGBA(p.cpicture, (*C.uint8_t)(cb), C.int(p.Stride))

	if err != nil {
		return err
	}

	if status != 1 {
		return fmt.Errorf("error importing RGBA: %d", status)
	}

	status, err = C.WebPEncode(conf, p.cpicture)

	if err != nil {
		return err
	}

	if status != 1 {
		return fmt.Errorf("error encoding RGBA: %d", status)
	}
	return nil
}
