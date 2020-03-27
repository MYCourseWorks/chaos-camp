package main

import (
	"fmt"
	"github.com/homework/hw2/cbitmap"
	"github.com/homework/hw2/path"
	"github.com/homework/hw2/vertex"
	"image/color"
	"os"
)

var bitmap *cbitmap.CBitMap
var white, black, blue, grayishBlue color.Color

func triangleExample() {
	p := path.New(blue, grayishBlue)
	p.PushBack(&vertex.Vertex{X: 400, Y: 400})
	p.PushBack(&vertex.Vertex{X: 300, Y: 400})
	p.PushBack(&vertex.Vertex{X: 300, Y: 300})
	p.PushBack(&vertex.Vertex{X: 300, Y: 200})
	p.PushBack(&vertex.Vertex{X: 400, Y: 400})

	p.Scale(50)
	// p.Rotate(-30)
	fmt.Println(p)

	bitmap.Clear(white)
	p.Draw(bitmap)
	file, _ := os.OpenFile("../out/image.png", os.O_RDWR|os.O_CREATE, os.ModePerm)
	bitmap.EncodePNG(file)
}

func complexPolygonExample() {
	p := path.New(black, blue)
	p.PushBack(&vertex.Vertex{X: 200, Y: 400})
	p.PushBack(&vertex.Vertex{X: 100, Y: 400})
	p.PushBack(&vertex.Vertex{X: 100, Y: 300})
	p.PushBack(&vertex.Vertex{X: 200, Y: 200})
	p.PushBack(&vertex.Vertex{X: 300, Y: 200})
	p.PushBack(&vertex.Vertex{X: 302, Y: 400})
	p.PushBack(&vertex.Vertex{X: 220, Y: 400})
	p.PushBack(&vertex.Vertex{X: 203, Y: 301})
	p.PushBack(&vertex.Vertex{X: 150, Y: 320})
	p.PushBack(&vertex.Vertex{X: 200, Y: 400})

	p.Scale(70)
	fmt.Println(p)

	bitmap.Clear(white)
	p.Draw(bitmap)
	file, _ := os.OpenFile("../out/image.png", os.O_RDWR|os.O_CREATE, os.ModePerm)
	bitmap.EncodePNG(file)
}

func init() {
	white = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	blue = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	grayishBlue = color.RGBA{R: 93, G: 133, B: 145, A: 255}
	bitmap = cbitmap.New(0, 0, 1028, 1028)
}

func main() {
	triangleExample()
}
