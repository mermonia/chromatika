package palette

import (
	"log"
	"testing"
)

func TestGeneratePalette(t *testing.T) {
	file := "test-1.jpg"
	_, err := GeneratePalette("/home/umbraslay/Pictures/Wallpapers/"+file, false)
	if err != nil {
		log.Fatal(err)
	}
}
