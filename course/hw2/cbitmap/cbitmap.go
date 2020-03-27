package cbitmap

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// CBitMap stands for Custom Bitmap
type CBitMap struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

// New creates a CBitMap instance
func New(x0, y0, x1, y1 int) *CBitMap {
	width, height := x1-x0, y1-y0
	pixels := make([]uint8, 4*width*height)
	return &CBitMap{
		Pix:    pixels,
		Stride: 4 * width,
		Rect:   image.Rect(x0, y0, x1, y1),
	}
}

// ColorModel retuns RGBAModel always
func (b *CBitMap) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the bound rectangle
func (b *CBitMap) Bounds() image.Rectangle {
	return b.Rect
}

// At returns the color of the pixel at (x, y).
func (b *CBitMap) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(b.Rect)) {
		return color.RGBA{}
	}
	i := b.PixOffset(x, y)
	s := b.Pix[i : i+4 : i+4]
	return color.RGBA{s[0], s[1], s[2], s[3]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (b *CBitMap) PixOffset(x, y int) int {
	return (y-b.Rect.Min.Y)*b.Stride + (x-b.Rect.Min.X)*4
}

// Set sets a pixel at x and y
func (b *CBitMap) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(b.Rect)) {
		return
	}
	i := b.PixOffset(x, y)
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	s := b.Pix[i : i+4 : i+4]
	s[0] = c1.R
	s[1] = c1.G
	s[2] = c1.B
	s[3] = c1.A
}

// Clear clears the bitmap
func (b *CBitMap) Clear(c color.Color) {
	for x := 0; x < b.Rect.Dx(); x++ {
		for y := 0; y < b.Rect.Dy(); y++ {
			b.Set(x, y, c)
		}
	}
}

// EncodePNG writes the image to a file
func (b *CBitMap) EncodePNG(file *os.File) error {
	err := png.Encode(file, b)
	return err
}
