package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/mermonia/chromatika/internal/clustering"
	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/evaluation"
	"github.com/mermonia/chromatika/internal/input"
	"github.com/mermonia/chromatika/internal/quantization"
)

type PaletteCommandOptions struct {
	ImagePath     string
	Clusters      int
	Fuzziness     float64
	QuantInterval int
	MaxIter       int
	Threshold     float64
	ScaleWidth    int
}

var paletteCommandDescription string = `
description goes here
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
	Commands: []*cli.Command{},
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
		cmdCfg := &PaletteCommandOptions{
			ImagePath:     filepath.Clean(c.StringArg("imagePath")),
			Clusters:      c.Int("clusters"),
			Fuzziness:     c.Float64("fuzziness"),
			QuantInterval: c.Int("quantInterval"),
			MaxIter:       c.Int("maxIter"),
			Threshold:     c.Float("threshold"),
			ScaleWidth:    c.Int("newWidth"),
		}

		return ExecutePalette(cmdCfg)
	},
}

func ExecutePalette(cmdCfg *PaletteCommandOptions) error {
	dominantColors, partMatrix, err := getDominantColors(cmdCfg)
	if err != nil {
		return fmt.Errorf("could not get dominant colors: %w", err)
	}

	primaryColor := evaluation.ChoosePrimaryColor(dominantColors, partMatrix)
	backgroundColor := evaluation.ChooseBackgroundColor(dominantColors, partMatrix)

	return nil
}

func getDominantColors(cmdCfg *PaletteCommandOptions) ([]*colors.Lab, [][]float64, error) {
	pixels, err := input.ReadImage(cmdCfg.ImagePath, cmdCfg.ScaleWidth)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read image: %w", err)
	}

	labColors, err := quantization.GetRawLabArray(pixels)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get lab color array: %w", err)
	}

	quantizedColors, np := quantization.Quantize(cmdCfg.QuantInterval, labColors)

	extractedColors, partMatrix, err := clustering.FCM(
		quantizedColors,
		np,
		cmdCfg.Fuzziness,
		cmdCfg.Threshold,
		cmdCfg.MaxIter,
		cmdCfg.Clusters,
	)

	return extractedColors, partMatrix, nil
}
