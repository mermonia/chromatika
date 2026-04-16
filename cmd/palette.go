package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/mermonia/chromatika/internal/format"
	"github.com/mermonia/chromatika/internal/palette"
	"github.com/urfave/cli/v3"
)

type PaletteCommandOptions struct {
	ImagePath     string
	Format        string
	Fuzziness     float64
	QuantInterval int
	MaxIter       int
	Threshold     float64
	ScaleWidth    int
	DarkMode      bool
}

var paletteCommandDescription string = `
Generates a complete, structured color palette from the given image.

The palette is built by repeatedly extracting dominant colors using
FCM clustering with varying parameters, selecting colors that
are sufficiently distinct. From this base set, additional colors are
derived to form a coherent palette, including light and dark neutrals,
accent colors, and harmonized variants.

Unlike 'extract', which returns raw dominant colors, this command
produces a full palette suitable for theming and UI usage.

Most stages of the extraction process can be tuned via flags, including
clustering behavior, quantization, and convergence thresholds. Image
downscaling can also be applied to improve performance.
`

var PaletteCommand cli.Command = cli.Command{
	Name:                  "palette",
	Aliases:               []string{"p"},
	Usage:                 "generate a palette from a background image",
	ArgsUsage:             "<path>",
	Description:           paletteCommandDescription,
	EnableShellCompletion: true,
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:  "imagePath",
			Value: "",
		},
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "quantInterval",
			Value:   4,
			Aliases: []string{"q"},
			Usage:   "adjust the color quantization interval",
		},
		&cli.IntFlag{
			Name:    "maxIter",
			Value:   100,
			Aliases: []string{"B"},
			Usage:   "adjust the maximum number of FCM iterations",
		},
		&cli.Float64Flag{
			Name:    "fuzziness",
			Value:   2.0,
			Aliases: []string{"m"},
			Usage:   "adjust the fuzziness of the algorithm's clusters",
		},
		&cli.Float64Flag{
			Name:    "threshold",
			Value:   0.001,
			Aliases: []string{"e"},
			Usage:   "adjust the min difference stop condition for fcm",
		},
		&cli.IntFlag{
			Name:    "newWidth",
			Value:   0,
			Aliases: []string{"w"},
			Usage:   "downscale the image proportionally to the given width",
		},
		&cli.BoolFlag{
			Name:    "darkmode",
			Value:   false,
			Aliases: []string{"d"},
			Usage:   "generate a darkmode palette",
		},
		&cli.StringFlag{
			Name:    "format",
			Value:   "ascii",
			Aliases: []string{"f"},
			Usage:   "output the generated palette in a specific format",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		cmdCfg := &PaletteCommandOptions{
			ImagePath:     filepath.Clean(c.StringArg("imagePath")),
			Fuzziness:     c.Float64("fuzziness"),
			QuantInterval: c.Int("quantInterval"),
			MaxIter:       c.Int("maxIter"),
			Threshold:     c.Float("threshold"),
			ScaleWidth:    c.Int("newWidth"),
			DarkMode:      c.Bool("darkmode"),
			Format:        c.String("format"),
		}

		return ExecutePalette(cmdCfg)
	},
}

func ExecutePalette(cmdCfg *PaletteCommandOptions) error {
	pal, err := palette.GeneratePalette(&palette.GenerationConfig{
		ImagePath:     cmdCfg.ImagePath,
		DarkMode:      cmdCfg.DarkMode,
		QuantInterval: cmdCfg.QuantInterval,
		ScaleWidth:    cmdCfg.ScaleWidth,
		Fuzziness:     cmdCfg.Fuzziness,
		Threshold:     cmdCfg.Threshold,
		MaxIter:       cmdCfg.MaxIter,
	})
	if err != nil {
		return fmt.Errorf("could not generate palette: %w", err)
	}

	var formatter format.PaletteFormatter

	switch cmdCfg.Format {
	case "ascii":
		formatter = &format.ASCIIPaletteFormatter{}
	case "toml":
		formatter = &format.TOMLPaletteFormatter{}
	default:
		return fmt.Errorf("unknown format specifier")
	}

	fmt.Print(formatter.Format(pal))

	return nil
}
