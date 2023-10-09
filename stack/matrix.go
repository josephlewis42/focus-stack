package stack

import (
	"image"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Float | constraints.Integer
}

type Matrix[T Number] struct {
	Values []T
	Width  int
	Height int
}

func NewMatrix[T Number](x int, y int) *Matrix[T] {
	return &Matrix[T]{
		Values: make([]T, x*y),
		Width:  x,
		Height: y,
	}
}

func (m *Matrix[T]) At(x, y int) T {
	return m.Values[m.offset(x, y)]
}

func (m *Matrix[T]) offset(x, y int) int {
	return y*m.Width + x
}

func (m *Matrix[T]) Set(x, y int, v T) {
	m.Values[m.offset(x, y)] = v
}

func (m *Matrix[T]) Each(callback func(x, y int, v T)) {
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			callback(x, y, m.At(x, y))
		}
	}
}

func (m *Matrix[T]) Bounds() image.Rectangle {
	return image.Rect(0, 0, m.Width-1, m.Height-1)
}

// MatrixFromPlane creates a matrix from a plane using the projection function converter.
func MatrixFromPlane[I any, O Number](input Plane[I], converter func(I) O) *Matrix[O] {
	bounds := input.Bounds().Canon()

	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	projected := NewMatrix[O](width, height)

	minX := bounds.Min.X
	minY := bounds.Min.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			projected.Set(x, y, converter(input.At(x+minX, y+minY)))
		}
	}

	return projected
}
