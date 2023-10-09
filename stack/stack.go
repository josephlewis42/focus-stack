package stack

import (
	"image"
	"image/color"
	"math"
)

const (
	// Width of the LoG Kernel to be applied before focus stacking.
	LoGKernelWidth = 13

	// Sigma of the Log Kernel to be applied before focus stacking.
	LogKernelSigma = 1.6
)

// FocusStack creates a stack and depth map from a set of images.
type FocusStack struct {
	kernel *Matrix[float32]

	// Number of images applied so far.
	numImages int

	// Maximum LoG seen at every position so far.
	maxLog *Matrix[float32]

	// Source image index for where each of the pixels came from.
	sourceImage *Matrix[int]

	// Output pixels, alpha channel is ignored.
	stackedImage *image.RGBA
}

func NewFocusStack(rect image.Rectangle) (*FocusStack, error) {
	kernel := LoGKernel(LogKernelSigma, LoGKernelWidth)

	return &FocusStack{
		kernel:       kernel,
		numImages:    0,
		maxLog:       NewMatrix[float32](rect.Max.X+1, rect.Max.Y+1),
		sourceImage:  NewMatrix[int](rect.Max.X+1, rect.Max.Y+1),
		stackedImage: image.NewRGBA(rect),
	}, nil
}

func grayscale(pixel color.Color) uint32 {
	r, g, b, _ := pixel.RGBA()
	// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
	// as those given by the JFIF specification and used by func RGBToYCbCr in
	// ycbcr.go.
	//
	// Note that 19595 + 38470 + 7471 equals 65536.
	//
	// The 24 is 16 + 8. The 16 is the same as used in RGBToYCbCr. The 8 is
	// because the return value is 8 bit color, not 16 bit color.
	return (19595*r + 38470*g + 7471*b + 1<<15) >> 24
}

func (stack *FocusStack) AddImage(inputImage image.Image) {
	stack.numImages++

	grayscaleImage := MatrixFromPlane[color.Color](inputImage, grayscale)
	filtered := NewFilteredImage[float32, uint32](stack.kernel, grayscaleImage)
	union := inputImage.Bounds().Union(stack.stackedImage.Rect)

	for y := union.Min.Y; y <= union.Max.Y; y++ {
		for x := union.Min.X; x <= union.Max.X; x++ {
			sharpness := filtered.At(x, y)

			sharper := math.Abs(float64(sharpness)) >= math.Abs(float64(stack.maxLog.At(x, y)))
			noData := stack.sourceImage.At(x, y) == 0
			if sharper || noData {
				stack.sourceImage.Set(x, y, stack.numImages)
				stack.maxLog.Set(x, y, sharpness)
				stack.stackedImage.Set(x, y, inputImage.At(x, y))
			}
		}
	}
}

// DepthLUT is a lookup table between an image index and color.
// Black is nearer to the camera.
type DepthLUT map[int]color.Gray16

// Gives a strict depthmap back to front assuming images were processed in order.
func (stack *FocusStack) OrderedDepths() DepthLUT {
	if stack.numImages == 0 {
		return nil
	}

	out := make(map[int]color.Gray16)

	step := float64(math.MaxUint16) / float64(stack.numImages)

	for i := 1; i <= stack.numImages; i++ {
		out[i] = color.Gray16{
			Y: math.MaxUint16 - uint16(math.Floor(float64(i)*step)),
		}
	}

	// Clamp first (implicit image) and last to max
	out[0] = color.Gray16{Y: 0}
	out[stack.numImages] = color.Gray16{Y: 255}

	return out
}

func (stack *FocusStack) DepthMap(depthLUT DepthLUT) *image.Gray16 {
	out := image.NewGray16(stack.stackedImage.Rect)

	stack.sourceImage.Each(func(x, y int, src int) {
		out.SetGray16(x, y, depthLUT[src])
	})

	return out
}

func (stack *FocusStack) StackedImage() image.Image {
	return stack.stackedImage
}
