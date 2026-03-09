package templating

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/mermonia/chromatika/internal/palette"
)

func RenderFileToWriter(src string, palette *palette.Palette, out io.Writer) error {
	t, err := template.ParseFiles(src)
	if err != nil {
		return fmt.Errorf("could not parse file for templating: %w", err)
	}
	return t.Execute(out, palette)
}

func RenderFileToFile(src, dst string, palette *palette.Palette) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("could not create parent dirs for dst: %w", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create dst file: %w", err)
	}
	defer out.Close()

	if err := RenderFileToWriter(src, palette, out); err != nil {
		return fmt.Errorf("could not render source file: %w", err)
	}

	return nil
}
