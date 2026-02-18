package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/mermonia/chromatika/internal/clustering"
	"github.com/mermonia/chromatika/internal/input"
	"github.com/mermonia/chromatika/internal/quantization"
)

type ExtractCommandOptions struct {
	ImagePath     string
	Clusters      int
	Fuzziness     float64
	QuantInterval int
	MaxIter       int
	Threshold     float64
	ScaleWidth    int
}

var extractCommandDescription string = `
description goes here
`

var ExtractCommand cli.Command = cli.Command{
	Name:                  "extract",
	Aliases:               []string{"e"},
	Usage:                 "extract dominant colors from an image",
	ArgsUsage:             "<path>",
	Description:           extractCommandDescription,
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
		cmdCfg := &ExtractCommandOptions{
			ImagePath:     filepath.Clean(c.StringArg("imagePath")),
			Clusters:      c.Int("clusters"),
			Fuzziness:     c.Float64("fuzziness"),
			QuantInterval: c.Int("quantInterval"),
			MaxIter:       c.Int("maxIter"),
			Threshold:     c.Float("threshold"),
			ScaleWidth:    c.Int("newWidth"),
		}

		return ExecuteExtract(cmdCfg)
	},
}

func ExecuteExtract(cmdCfg *ExtractCommandOptions) error {
	pixels, err := input.ReadImage(cmdCfg.ImagePath, cmdCfg.ScaleWidth)
	if err != nil {
		return fmt.Errorf("could not read image: %w", err)
	}

	labColors, err := quantization.GetRawLabArray(pixels)
	if err != nil {
		return fmt.Errorf("could not get lab color array: %w", err)
	}

	quantizedColors, np := quantization.Quantize(cmdCfg.QuantInterval, labColors)

	extractedColors, _, err := clustering.FCM(
		quantizedColors,
		np,
		cmdCfg.Fuzziness,
		cmdCfg.Threshold,
		cmdCfg.MaxIter,
		cmdCfg.Clusters,
	)

	for _, color := range extractedColors {
		render, err := color.Render(3)
		if err != nil {
			return fmt.Errorf("could not render color: %w", err)
		}
		fmt.Print(render)
	}
	fmt.Println()

	return nil
}
