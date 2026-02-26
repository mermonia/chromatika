package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/mermonia/chromatika/internal/clustering"
	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/evaluation"
	"github.com/mermonia/chromatika/internal/extraction"
	"github.com/mermonia/chromatika/internal/palette"
	"github.com/urfave/cli/v3"
)

type PaletteImageCommandOptions struct {
	ImagePath     string
	Clusters      int
	Fuzziness     float64
	QuantInterval int
	MaxIter       int
	Threshold     float64
	ScaleWidth    int
}

var paletteImageCommandDescripton string = `
desc goes here
`

var PaletteImageCommand cli.Command = cli.Command{
	Name:                  "image",
	Aliases:               []string{"i"},
	Usage:                 "generate palette from a background image",
	ArgsUsage:             "<path>",
	Description:           paletteImageCommandDescripton,
	EnableShellCompletion: true,
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:  "imagePath",
			Value: "",
		},
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "clusters",
			Value:   8,
			Aliases: []string{"k"},
			Usage:   "adjust the number of colors extracted",
		},
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
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		cmdCfg := &PaletteImageCommandOptions{
			ImagePath:     filepath.Clean(c.StringArg("imagePath")),
			Clusters:      c.Int("clusters"),
			Fuzziness:     c.Float64("fuzziness"),
			QuantInterval: c.Int("quantInterval"),
			MaxIter:       c.Int("maxIter"),
			Threshold:     c.Float("threshold"),
			ScaleWidth:    c.Int("newWidth"),
		}

		return ExecutePaletteImage(cmdCfg)
	},
}

func ExecutePaletteImage(cmdCfg *PaletteImageCommandOptions) error {
	dominantColors, partMatrix, err := extraction.GetDominantColors(
		cmdCfg.ImagePath,
		cmdCfg.ScaleWidth,
		cmdCfg.QuantInterval,
		clustering.FCMParameters{
			M: cmdCfg.Fuzziness,
			E: cmdCfg.Threshold,
			B: cmdCfg.MaxIter,
			K: cmdCfg.Clusters,
		},
	)

	if err != nil {
		return fmt.Errorf("could not get dominant colors: %w", err)
	}

	primaryColor := evaluation.ChoosePrimaryColor(dominantColors, partMatrix)
	backgroundColor := evaluation.ChooseBackgroundColor(dominantColors, partMatrix)

	primaryColorBlock, _ := primaryColor.Render(3)
	backgroundColorBlock, _ := backgroundColor.Render(3)

	fmt.Printf("Primary color: %s\n", primaryColorBlock)
	fmt.Printf("Background color: %s\n", backgroundColorBlock)

	darkmode := backgroundColor.L < 50

	neutrals := palette.GenerateNeutrals(colors.LabToLCH(backgroundColor), darkmode)
	fmt.Print(neutrals)

	accents := palette.GenerateAccents(colors.LabToLCH(primaryColor), darkmode)
	fmt.Print(accents)

	return nil
}
