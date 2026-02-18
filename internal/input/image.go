package input

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

func ReadImage(path string, newW int) ([]color.Color, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	defer file.Close()

	pixels, err := getPixels(file, newW)
	if err != nil {
		return nil, fmt.Errorf("could not read the file's pixels: %w", err)
	}

	return pixels, err
}

func getPixels(r io.Reader, newW int) ([]color.Color, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("could not decode given file: %w", err)
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	if newW != 0 && newW < width {
		scaleFactor := width / newW
		width, height = width/scaleFactor, height/scaleFactor

		img = scaleImage(img, width, height)
		bounds = img.Bounds()
	}

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

func scaleImage(src image.Image, newW, newH int) image.Image {
	var dst draw.Image = image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}
