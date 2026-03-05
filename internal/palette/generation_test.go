package palette

import (
	"fmt"
	"log"
	"testing"
)

func TestGeneratePalette(t *testing.T) {
	file := "test-9.jpg"
	_, err := GeneratePalette("/home/umbraslay/Pictures/Wallpapers/"+file, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("palette end")
}
