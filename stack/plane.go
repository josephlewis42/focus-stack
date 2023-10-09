package stack

import "image"

// Plane represents a 2D plane.
type Plane[T any] interface {
	At(x, y int) T
	Bounds() image.Rectangle
}
