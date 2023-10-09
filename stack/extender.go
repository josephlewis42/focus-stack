package stack

import (
	"image"
	"math"
)

// Extender extends the image pixel space indefintely by picking
// the nearest pixel to the point.
type Extender[T any] struct {
	source Plane[T]
	bounds image.Rectangle
}

func NewExtender[T any](source Plane[T]) *Extender[T] {
	return &Extender[T]{
		source: source,
		bounds: source.Bounds().Canon(),
	}
}

func (e *Extender[T]) At(x, y int) T {
	return e.source.At(
		e.nearestX(x),
		e.nearestY(y),
	)
}

func (e *Extender[T]) nearestX(x int) int {
	switch {
	case x < e.bounds.Min.X:
		return e.bounds.Min.X
	case x > e.bounds.Max.X:
		return e.bounds.Max.X
	default:
		return x
	}
}

func (e *Extender[T]) nearestY(y int) int {
	switch {
	case y < e.bounds.Min.Y:
		return e.bounds.Min.Y
	case y > e.bounds.Max.Y:
		return e.bounds.Max.Y
	default:
		return y
	}
}

func (e *Extender[T]) Bounds() image.Rectangle {
	return image.Rect(math.MinInt, math.MinInt, math.MaxInt, math.MaxInt)
}
