package stack

import (
	"math"
)

// Compute the Laplacian of Gaussian for a given point.
func LoG(sigma float64, x int, y int) float64 {
	// https://homepages.inf.ed.ac.uk/rbf/HIPR2/log.htm
	xSq := float64(x * x)
	ySq := float64(y * y)
	sSq := float64(sigma * sigma)

	return (-1 / (math.Pi * math.Pow(sigma, 4))) *
		(1 - (xSq+ySq)/(2*sSq)) *
		math.Exp(-(xSq+ySq)/(2*sSq))
}

// Computes a Laplacian of Gaussian kernel of size N (or n+1 if n is even).
func LoGKernel(sigma float64, n int) *Matrix[float32] {
	if n%2 == 0 {
		n = n + 1
	}

	kernel := NewMatrix[float32](n, n)

	midpoint := (n - 1) / 2
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			// The output is symmetric so we could be more efficient here
			// but it doesn't really matter for small kernels.
			kernel.Set(i, j, float32(LoG(sigma, i-midpoint, j-midpoint)))
		}
	}

	return kernel
}
