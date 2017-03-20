package img2ansi

import (
	"fmt"
	"image"
	"image/color"
	"io"
)

const (
	ansiColorBase  int     = 16
	ansiColorSteps float64 = 6
	rgbaColorSpace float64 = float64(1 << 16)
)

// ImageToAnsi writes the converted ANSI image to the writer
func ImageToAnsi(img image.Image, w io.Writer) (err error) {
	var previous string
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			current := toAnsiCode(img.At(x, y))
			if current != previous {
				_, err = fmt.Fprint(w, current)
				if err != nil {
					return err
				}
			}
			_, err = fmt.Fprint(w, "0")
			if err != nil {
				return err
			}
			previous = current
		}
		_, err = fmt.Fprintln(w)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprint(w, "\x1b[0m")
	if err != nil {
		return err
	}
	return nil
}

func toAnsiCode(c color.Color) string {
	r, g, b, _ := c.RGBA()
	// note: there are 256 colors in total in the xterm color space
	//       the first 16 colors are 8 base colors and 8 bright colors (skipped with `ansiColorBase`)
	//       the next 216 colors are all colors in steps of 6
	//       the last 24 are a grayscale from black to white
	code := ansiColorBase + toAnsiSpace(r)*36 + toAnsiSpace(g)*6 + toAnsiSpace(b)
	return fmt.Sprintf("\x1b[38;5;%dm", code)
}

func toAnsiSpace(val uint32) int {
	return int(ansiColorSteps * (float64(val) / rgbaColorSpace))
}
