package quantization

import (
	"fmt"
	"image/color"
	"runtime"
	"sync"

	"github.com/mermonia/chromatika/internal/colors"
)

func GetRawLabArray(pixels []color.Color) ([]*colors.Lab, error) {
	workers := runtime.NumCPU()

	res := make([]*colors.Lab, len(pixels))

	var wg sync.WaitGroup
	wg.Add(workers)

	errCh := make(chan error, workers)

	chunkSize := (len(pixels) + workers - 1) / workers

	for w := range workers {
		start := w * chunkSize
		end := start + chunkSize
		end = min(end, len(pixels))

		go func(start, end int) {
			defer wg.Done()

			for i := start; i < end; i++ {
				lab, err := colors.PixelToLab(pixels[i])
				if err != nil {
					errCh <- fmt.Errorf("could not convert pixel to lab: %w", err)
					return
				}
				res[i] = lab
			}
		}(start, end)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
