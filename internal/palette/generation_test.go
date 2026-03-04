package palette

import (
	"fmt"
	"log"
	"testing"
)

func TestGeneratePalette(t *testing.T) {
	_, err := GeneratePalette("/home/umbraslay/Pictures/Wallpapers/hatsune-miku-dark-2k.png", false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("palette end")
}
