package path

import (
	"container/list"
	"fmt"
	"github.com/homework/hw2/cbitmap"
	"github.com/homework/hw2/vertex"
	"image"
	"image/color"
)

// Path is a structure that represents a set of points in a 2D plane.
type Path struct {
	path        *list.List
	StrokeColor color.Color // TODO: all style options could be grouped together in another structure
	FillColor   color.Color
}

// New creates an instance of the Path structure.
func New(stroke color.Color, fillColor color.Color) *Path {
	p := &Path{
		path:        list.New(),
		StrokeColor: stroke,
		FillColor:   fillColor,
	}

	return p
}

// ToArray converts the path to an array
func (p *Path) ToArray() []vertex.Vertex {
	ret := make([]vertex.Vertex, p.path.Len())
	p.Iterate(func(l *list.Element, i int) bool {
		v := l.Value.(*vertex.Vertex)
		ret[i] = vertex.Vertex{X: v.X, Y: v.Y}
		return false
	})

	return ret
}

// Iterate traverses every element in the path.
// NOTE: this could be an extension on the list.List struct for more generic use!
func (p *Path) Iterate(cb func(*list.Element, int) bool) {
	i := 0
	currElem := p.path.Front()
	for currElem != nil {
		breakFlag := cb(currElem, i)
		if breakFlag {
			return
		}
		currElem = currElem.Next()
		i++
	}
}

// Len returns the number of points in the path.
func (p *Path) Len() int {
	return p.path.Len()
}

// Distance calculate the overall vector distance in the path.
func (p *Path) Distance() (dist float64) {
	if p.path == nil || p.path.Len() == 0 {
		return 0
	}

	dist = 0
	currElem := p.path.Front()
	for currElem != nil && currElem.Next() != nil {
		currVector := currElem.Value.(*vertex.Vertex)
		nextVector := currElem.Next().Value.(*vertex.Vertex)
		dist += currVector.Distance(nextVector)
		currElem = currElem.Next()
	}

	return dist
}

// Translate moves every point in the direction of the vector parameter.
func (p *Path) Translate(vector *vertex.Vertex) {
	p.Iterate(func(l *list.Element, i int) bool {
		v := l.Value.(*vertex.Vertex)
		v.Translate(vector)
		return false
	})
}

// Scale scales a path by a given factor.
// Uses Coordinate normalization -
// Mapping entries of a data set from their natural range to values between 0 and 1.
// And translates them by a given factor.
func (p *Path) Scale(f float64) {
	maxX, minX, minY, maxY := p.findMinAndMax()
	if maxX-minX == 0 || maxY-minY == 0 {
		panic("Dividing by zero")
	}

	T := vertex.Vertex{X: -minX, Y: -minY}
	Sx := 1 / (maxX - minX)
	Sy := 1 / (maxY - minY)

	p.Iterate(func(l *list.Element, i int) bool {
		v := l.Value.(*vertex.Vertex)
		v.X += ((v.X + T.X) * Sx) * f
		v.Y += ((v.Y + T.Y) * Sy) * f
		return false
	})
}

// TODO: we can memorise minX, maxX, minY, maxY ?
func (p *Path) findMinAndMax() (minX, maxX, minY, maxY float64) {
	if p.Len() <= 0 {
		return 0, 0, 0, 0
	}

	maxX = p.path.Front().Value.(*vertex.Vertex).X
	minX = p.path.Front().Value.(*vertex.Vertex).X
	minY = p.path.Front().Value.(*vertex.Vertex).Y
	maxY = p.path.Front().Value.(*vertex.Vertex).Y

	p.Iterate(func(l *list.Element, i int) bool {
		v := l.Value.(*vertex.Vertex)
		switch {
		case v.X > maxX:
			maxX = v.X
		case v.X < minX:
			minX = v.X
		case v.Y > maxY:
			maxY = v.Y
		case v.Y < minY:
			minY = v.Y
		}

		return false
	})

	return maxX, minX, minY, maxY
}

// Rotate rotates the path by a given angle in degrees.
func (p *Path) Rotate(angle float64) {
	if p.Len() <= 0 {
		return
	}

	prevPosition := p.path.Front().Value.(*vertex.Vertex)
	currElem := p.path.Front()
	for currElem != nil && currElem.Next() != nil {
		currVector := currElem.Value.(*vertex.Vertex)
		nextVector := currElem.Next().Value.(*vertex.Vertex)
		currElem = currElem.Next()

		d := vertex.Vertex{
			X: nextVector.X - prevPosition.X,
			Y: nextVector.Y - prevPosition.Y,
		}

		prevPosition = &vertex.Vertex{X: nextVector.X, Y: nextVector.Y}
		nextVector.X = currVector.X
		nextVector.Y = currVector.Y
		nextVector.Translate(&d)
		nextVector.Rotate(currVector, angle)
		fmt.Println()
	}
}

// Add adds a vector to the path
// TODO: how do we handle out of range ?
func (p *Path) Add(vector *vertex.Vertex, position int) {
	p.Iterate(func(l *list.Element, i int) bool {
		v := l.Value.(*vertex.Vertex)
		if i == position {
			p.path.InsertAfter(v, l)
			return true
		}

		return false
	})
}

// PushBack pushes a vector at the end of the path
func (p *Path) PushBack(vector *vertex.Vertex) {
	p.path.PushBack(vector)
}

// Remove removes a vector at a given position
// TODO: how do we handle out of range ?
func (p *Path) Remove(position int) {
	p.Iterate(func(l *list.Element, i int) bool {
		if i == position {
			p.path.Remove(l)
			return true
		}

		return false
	})
}

// Draw draws the path. If the path forms a polygon, it fills it.
// TODO: this
func (p *Path) Draw(bitmap *cbitmap.CBitMap) image.Image {
	p.stroke(bitmap)

	// FIXME: this code should be replaced with a
	// proper implementation of the scan fill algorithm.
	if p.path.Len() > 2 {
		front := p.path.Front().Value.(*vertex.Vertex)
		back := p.path.Back().Value.(*vertex.Vertex)

		// Naive check if the path is a Polygon :
		if front.IsEqualTo(back) {

			width := bitmap.Rect.Dx()
			height := bitmap.Rect.Dy()

			// TODO:
			// if min and max for X and Y are saved,
			// we could traverse a sub matrix instead
			// of the whole bitmap image.
			for y := 0; y < width; y++ {

				// Scan the entire line and find all points that are part of the stroke.
				ray := make([]int, 0)
				for x := 0; x < height; x++ {
					r, g, b, a := bitmap.At(x, y).RGBA()
					sr, sg, sb, sa := p.StrokeColor.RGBA()
					if r == sr && g == sg && b == sb && a == sa {
						ray = append(ray, x)
					}
				}

				rayLen := len(ray)
				isInPolygon := false
				for x := 0; x < rayLen-1; x++ {
					if ray[x] != ray[x+1]-1 {
						isInPolygon = !isInPolygon
					}

					if isInPolygon {
						for i := ray[x] + 1; i < ray[x+1]; i++ {
							bitmap.Set(i, y, p.FillColor)
						}
					}
				}
			}
		}
	}

	return bitmap
}

func (p *Path) stroke(bitmap *cbitmap.CBitMap) {
	currElem := p.path.Front()
	for currElem != nil && currElem.Next() != nil {
		currVector := currElem.Value.(*vertex.Vertex)
		nextVector := currElem.Next().Value.(*vertex.Vertex)
		currElem = currElem.Next()

		drawLine(
			bitmap,
			int(currVector.X), // TODO: should round instead of droping the fractional bits
			int(currVector.Y),
			int(nextVector.X),
			int(nextVector.Y),
			&p.StrokeColor,
		)
	}
}

// Bresenham's line algorithm
func drawLine(img *cbitmap.CBitMap, x0, y0, x1, y1 int, color *color.Color) {
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}

	sx := 0
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	sy := 0
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		img.Set(x0, y0, *color)

		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}
