package colors

import (
	"fmt"
	"log"
	"testing"
)

func TestRGBtoXYZ(t *testing.T) {
	a := &Rgb{
		R: 200,
		G: 0,
		B: 50,
	}
	xyz, err := RGBtoXYZ(a)
	if err != nil {
		log.Fatalf("could not convert: %s", err.Error())
	}
	fmt.Println(xyz.X, xyz.Y, xyz.Z)
}

func TestRGBtoLab(t *testing.T) {
	a := &Rgb{
		R: 255,
		G: 255,
		B: 255,
	}

	_, err := RGBtoLab(a)
	if err != nil {
		log.Fatalf("could not convert: %s", err.Error())
	}
}
