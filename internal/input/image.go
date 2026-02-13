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
	minX, minY := img.Bounds().Min.X, img.Bounds().Min.Y
	maxX, maxY := img.Bounds().Max.X, img.Bounds().Max.Y

	pixels := make([]color.Color, bounds.Dx()*bounds.Dy())

	for i := minX; i < maxX; i++ {
		for j := minY; j < maxY; j++ {
			idx := i*bounds.Dx() + j
			pixels[idx] = img.At(i, j)
		}
	}

	return pixels, nil
}
