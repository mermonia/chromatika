package input

import (
	"fmt"
	_ "golang.org/x/image/webp"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

func ReadImage(path string) ([]color.Color, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	defer file.Close()

	pixels, err := getPixels(file)
	if err != nil {
		return nil, fmt.Errorf("could not read the file's pixels: %w", err)
	}

	return pixels, err
}

func getPixels(r io.Reader) ([]color.Color, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("could not decode given file: %w", err)
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	pixels := make([]color.Color, width*height)

	idx := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixels[idx] = img.At(x, y)
			idx++
		}
	}

	return pixels, nil
}
