package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/lloiser/img2ansi"
	"github.com/nfnt/resize"
)

var outputPtr = flag.String("o", "", "Output file")
var widthPtr = flag.Int("w", 100, "Width of ANSI image")

// note: this is the width/height ratio of a typical
// fixed-width font (+-0.03 for other fonts)
const ratio = 0.5

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] <image>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		exit("No image provided")
	}

	imageName := args[0]
	imageFile, err := os.Open(imageName)
	if err != nil {
		exit("Could not open image %q", imageName)
	}
	defer imageFile.Close()

	imageReader := bufio.NewReader(imageFile)
	img, _, err := image.Decode(imageReader)
	if err != nil {
		exit("Could not decode image")
	}

	var writer io.Writer
	if outputPtr != nil && *outputPtr != "" {
		outputFile, err := os.Create(*outputPtr)
		if err != nil {
			exit("Could not open %q for writing", *outputPtr)
		}
		defer outputFile.Close()
		writer = outputFile
	} else {
		writer = os.Stdout
	}

	if widthPtr != nil && *widthPtr > 0 {
		img = resizeImg(img, *widthPtr)
	}

	img2ansi.ImageToAnsi(img, writer)
}

func exit(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(2)
}

func resizeImg(img image.Image, width int) image.Image {
	max := img.Bounds().Max
	curRatio := float64(max.Y) / float64(max.X)
	if width > max.X {
		width = max.X
	}
	w := uint(width)
	h := uint(float64(width) * curRatio * ratio)
	return resize.Resize(w, h, img, resize.Lanczos3)
}
