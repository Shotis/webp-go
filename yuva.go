package webp

import "image"

type YUVAPicture struct {
	bounds           image.Rectangle
	y, u, v          []uint8
	yStride          int
	uvStride         int
	subsamplingRatio image.YCbCrSubsampleRatio
}

func (p *YUVAPicture) YStride() int {
	return p.yStride
}

func (p *YUVAPicture) UVStride() int {
	return p.uvStride
}

func (*YUVAPicture) RGBA() bool {
	return false
}

func (p *YUVAPicture) Bounds() image.Rectangle {
	return p.bounds
}

func NewYUVAPicture(img *image.YCbCr) *YUVAPicture {
	return &YUVAPicture{
		bounds:           img.Bounds(),
		y:                img.Y,
		u:                img.Cb,
		v:                img.Cr,
		yStride:          img.YStride,
		uvStride:         img.CStride,
		subsamplingRatio: img.SubsampleRatio,
	}
}
