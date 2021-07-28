package webp

/*
#cgo pkg-config: libwebp
#include <stdlib.h>
#include "webp_util.h"
*/
import "C"
import (
	"fmt"
	"image"
	"io"
	"strings"
	"unsafe"
)

type Hint int

const (
	DefaultHint Hint = iota
	PictureHint
	PhotoHint
	GraphHint
	LastHint
)

var hints = map[string]Hint{
	"default": DefaultHint,
	"picture": PictureHint,
	"photo":   PhotoHint,
	"graph":   GraphHint,
}

func GetHint(s string) Hint {
	v, ok := hints[strings.ToLower(s)]

	if !ok {
		return DefaultHint
	}
	return v
}

type Picture interface {
	Bounds() image.Rectangle
	RGBA() bool
}

type allocatedPicture struct {
	Picture
	cpicture *C.WebPPicture
}

func NewPicture(img image.Image) *allocatedPicture {
	switch p := img.(type) {
	case *image.RGBA:
		return &allocatedPicture{Picture: NewRGBAImage(p)}
	case *image.NRGBA:
		return &allocatedPicture{Picture: NewNRGBAImage(p)}
	case *image.YCbCr:
		return &allocatedPicture{Picture: NewYUVAPicture(p)}
	}
	return nil
}

func (p *allocatedPicture) Init() error {
	var pic C.struct_WebPPicture
	i, err := C.WebPPictureInit(&pic)
	if err != nil {
		return err
	}
	if i != 1 {
		return fmt.Errorf("error occurred initializing the picutre: %d", i)
	}

	pic.width = C.int(p.Bounds().Dx())
	pic.height = C.int(p.Bounds().Dy())

	err = p.alloc()

	if err != nil {
		return err
	}

	p.cpicture = &pic

	return nil
}

func (p *allocatedPicture) alloc() error {
	i, err := C.WebPPictureAlloc(p.cpicture)

	if err != nil {
		return err
	}

	if i != 1 {
		return fmt.Errorf("error occurred allocating picture: %d", i)
	}
	return nil
}

func (p *allocatedPicture) Free() error {
	_, err := C.WebPPictureFree(p.cpicture)

	if err != nil {
		return err
	}
	return nil
}

func (p *allocatedPicture) importRGBA(pic *RGBAPicture) error {

	cb := C.CBytes([]byte(pic.Pixels())) // translate the pixels to
	defer C.free(cb)

	// import our pixels into the image
	status, err := C.WebPPictureImportRGBA(p.cpicture, (*C.uint8_t)(cb), C.int(pic.Stride()))

	if err != nil {
		return err
	}

	if status != 1 {
		return fmt.Errorf("error importing RGBA: %d", status)
	}
	return nil
}

func (p *allocatedPicture) importYUVA(pic *YUVAPicture) {
	// YUVA
	y := C.CBytes([]byte(pic.y))
	u := C.CBytes([]byte(pic.u))
	v := C.CBytes([]byte(pic.v))

	defer func() {
		C.free(y)
		C.free(u)
		C.free(v)
	}()

	p.cpicture.y = (*C.uint8_t)(y)
	p.cpicture.u = (*C.uint8_t)(u)
	p.cpicture.v = (*C.uint8_t)(v)

	p.cpicture.y_stride = C.int(pic.yStride)
	p.cpicture.uv_stride = C.int(pic.uvStride)
}

func (p *allocatedPicture) Encode(w io.Writer, config *Config) error {
	conf, err := config.toCStruct()

	if err != nil {
		return err
	}

	p.cpicture.writer = C.WebPWriterFunction(C.writeTo)
	wh.add(unsafe.Pointer(p.cpicture), w)

	defer wh.del(unsafe.Pointer(p.cpicture))

	if p.Picture.RGBA() {
		// encode RGBA
		if err := p.importRGBA(p.Picture.(*RGBAPicture)); err != nil {
			return err
		}
	} else {
		p.importYUVA(p.Picture.(*YUVAPicture))
	}

	status, err := C.WebPEncode(conf, p.cpicture)

	if err != nil {
		return err
	}

	if status != 1 {
		return fmt.Errorf("error encoding the image: %d", status)
	}
	return nil
}
