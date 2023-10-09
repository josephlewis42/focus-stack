package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"runtime/pprof"
	"sort"

	"github.com/josephlewis42/focus-stack/stack"

	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"flag"
)

func main() {

	outputPath := flag.String("o", "stacked.png", "Output file for stacked image.")
	depthPath := flag.String("d", "depth.png", "Output file for depth map.")
	cpuprofile := flag.String("cpuprofile", "", "Write cpu profile to file.")
	flag.Parse()

	files := flag.Args()
	sort.Strings(files)

	// Pre-open output files so we can bail early
	stackedFd, err := os.Create(*outputPath)
	if err != nil {
		log.Fatalf("couldn't open %q: %v", *outputPath, err)
	}
	defer stackedFd.Close()

	depthFd, err := os.Create(*depthPath)
	if err != nil {
		log.Fatalf("couldn't open %q: %v", *depthPath, err)
	}
	defer depthFd.Close()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Process images
	var images []image.Image
	for _, fileName := range files {
		fmt.Fprintln(os.Stderr, "Reading", fileName)
		fd, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("couldn't open image %q: %v\n", fileName, err)
		}

		img, _, err := image.Decode(fd)
		if err != nil {
			log.Fatalf("couldn't decode image %q: %v\n", fileName, err)
		}

		images = append(images, img)
	}

	fmt.Fprintln(os.Stderr, "Setting up stack")
	if len(images) == 0 {
		log.Fatalln("Require at least one image argument")
		return
	}

	focusStack, err := stack.NewFocusStack(images[0].Bounds())
	if err != nil {
		log.Fatalf("couldn't create focus stack: %v\n", err)
	}

	for idx, img := range images {
		fmt.Fprintln(os.Stderr, "Processing", idx)
		focusStack.AddImage(img)
	}

	png.Encode(stackedFd, focusStack.StackedImage())
	png.Encode(depthFd, focusStack.DepthMap(focusStack.OrderedDepths()))
}
