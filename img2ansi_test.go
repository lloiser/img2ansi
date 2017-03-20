package img2ansi_test

import (
	"bytes"
	"flag"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/lloiser/img2ansi"
)

func TestImageToAnsiSmall(t *testing.T) {
	img := fileToImg(t, "testdata/small.gif")
	buf := &bytes.Buffer{}
	err := img2ansi.ImageToAnsi(img, buf)
	if err != nil {
		t.Fatal(err)
	}
	expected := "\x1b[38;5;16m0000\n" + // 4x black
		"0\x1b[38;5;196m0\x1b[38;5;46m0\x1b[38;5;16m0\n" + // 1x black 1x red 1x green 1x black
		"0\x1b[38;5;21m0\x1b[38;5;231m0\x1b[38;5;16m0\n" + // 1x black 1x blue 1x white 1x black
		"0000\n" + // 4x black
		"\x1b[0m" // reset
	if buf.String() != expected {
		t.Fatalf("Expected %q but got %q", expected, buf)
	}
}

func fileToImg(t *testing.T, file string) image.Image {
	f, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	return img
}

var goldenPtr = flag.Bool("golden", false, "Update the testdata/*.golden files")

func TestVariousImages(t *testing.T) {
	flag.Parse()

	images := []string{
		"gopher.png",
		"smiley.gif",
	}
	goldenFile := func(file string) string {
		return "testdata/" + path.Base(file) + ".golden"
	}

	for _, file := range images {
		img := fileToImg(t, "testdata/"+file)
		actual := &bytes.Buffer{}
		err := img2ansi.ImageToAnsi(img, actual)
		if err != nil {
			t.Fatalf("%s: %s", file, err)
		}

		expectedFile := goldenFile(file)
		expected, _ := ioutil.ReadFile(expectedFile)
		if !bytes.Equal(actual.Bytes(), expected) {
			if goldenPtr != nil && *goldenPtr {
				err = ioutil.WriteFile(expectedFile, actual.Bytes(), 0666)
				if err != nil {
					t.Fatalf("%s: %s", file, err)
				}
				continue
			}
			t.Fatalf("%s: Expected %q but got %q", file, expected, actual)
		}
	}
}
