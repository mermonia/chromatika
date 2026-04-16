package cmd

import (
	"context"
	"fmt"

	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/input"
	"github.com/urfave/cli/v3"
)

type ColorType int

const (
	Rgb ColorType = iota
	Lab
)

type AnalyzeCommandOptions struct {
	ColorString string
	ColType     ColorType
}

var analyzeCommandDescription string = `
Prints the LCH(ab) color equivalent to the given sRGB color in hex
notation (#RRGGBB).

While this command is meant to be used as a debugging/testing tool for
colors generated with chromatiak, it can be used to analyze any RGB color 
in a more expressive color space.
`

var AnalyzeCommand cli.Command = cli.Command{
	Name:                  "analyze",
	Aliases:               []string{"a"},
	Usage:                 "analyze the given color",
	ArgsUsage:             "<color>",
	Description:           analyzeCommandDescription,
	EnableShellCompletion: true,
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:  "color",
			Value: "",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		var colorFormat ColorType = Rgb

		cmdCfg := &AnalyzeCommandOptions{
			ColorString: c.StringArg("color"),
			ColType:     colorFormat,
		}

		return ExecuteAnalyze(cmdCfg)
	},
}

func ExecuteAnalyze(cmdCfg *AnalyzeCommandOptions) error {
	rgb := input.RGBfromString(cmdCfg.ColorString)
	fmt.Printf("Color to analyze: %v\n", rgb)

	lab, err := colors.RGBtoLab(rgb)
	if err != nil {
		return fmt.Errorf("could not convert to lab: %w", err)
	}
	lch := colors.LabToLCH(lab)

	fmt.Printf("Luminosity: %f\n", lch.L)
	fmt.Printf("Chroma: %f\n", lch.C)
	fmt.Printf("Hue: %f\n", lch.H)

	return nil
}
