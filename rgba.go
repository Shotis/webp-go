package webp

import "image"

type RGBAPicture struct {
	bounds image.Rectangle
	pixels []uint8
	stride int
}

func (p *RGBAPicture) Bounds() image.Rectangle {
	return p.bounds
}

func (p *RGBAPicture) Pixels() []uint8 {
	return p.pixels
}

func (p *RGBAPicture) Stride() int {
	return p.stride
}

func (p *RGBAPicture) RGBA() bool {
	return true
}

func NewNRGBAImage(img *image.NRGBA) *RGBAPicture {
	return &RGBAPicture{
		bounds: img.Bounds(),
		pixels: img.Pix,
		stride: img.Stride,
	}
}

func NewRGBAImage(img *image.RGBA) *RGBAPicture {
	return &RGBAPicture{
		bounds: img.Bounds(),
		pixels: img.Pix,
		stride: img.Stride,
	}
}
