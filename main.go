package main

import (
	"fmt"
	"log"

	"github.com/mermonia/chromatika/internal/clustering"
	"github.com/mermonia/chromatika/internal/input"
	"github.com/mermonia/chromatika/internal/quantization"
)

func main() {
	cols, err := input.ReadImage("./test-image.png")
	if err != nil {
		log.Fatalf("could not read image: %s", err.Error())
	}

	labCols, err := quantization.GetRawLabArray(cols)
	if err != nil {
		log.Fatalf("could not get raw lab arrayl: %s", err.Error())
	}

	qColArray, np := quantization.Quantize(4, labCols)

	dominantColors, U, err := clustering.FCM(qColArray, np, 2.0, 0.001, 100, 8)
	if err != nil {
		log.Fatalf("could not execute FCM: %s", err.Error())
	}

	fmt.Println(len(dominantColors), len(U))
	for _, color := range dominantColors {
		fmt.Println(color.L, color.A, color.B)
	}
}
