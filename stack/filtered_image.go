package stack

import (
	"image"
)

// FilteredImage applies a filter to an extended image
type FilteredImage[K, E Number] struct {
	kernel   *Matrix[K]
	extended Plane[E]
	original Plane[E]
}

func NewFilteredImage[K, E Number](kernel *Matrix[K], backend Plane[E]) *FilteredImage[K, E] {
	return &FilteredImage[K, E]{
		kernel:   kernel,
		extended: NewExtender(backend),
		original: backend,
	}
}

func (img *FilteredImage[K, E]) At(x, y int) K {
	kernelMidX := (img.kernel.Width - 1) / 2
	kernelMidY := (img.kernel.Height - 1) / 2

	var value K
	for ky := 0; ky < img.kernel.Height; ky++ {
		for kx := 0; kx < img.kernel.Width; kx++ {
			kVal := img.kernel.At(kx, ky)

			pixel := img.extended.At(x+kx-kernelMidX, y+ky-kernelMidY)
			value += K(float32(kVal) * float32(pixel))
		}
	}

	return value
}

func (img *FilteredImage[K, E]) Bounds() image.Rectangle {
	return img.original.Bounds()
}
