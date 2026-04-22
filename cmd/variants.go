package cmd

import (
	"context"
	"fmt"

	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/format"
	"github.com/mermonia/chromatika/internal/input"
	"github.com/mermonia/chromatika/internal/variants"
	"github.com/urfave/cli/v3"
)

type VariantsCommandOptions struct {
	Color  string
	Format string
}

var variantsCommandDescription string = `Generate a set of variants from a given color.

Chromatika will try to make the given color as central as possible in 
the output set, while allowing room for all 10 luminosity variants (for
example, providing a color with a medium luminosity will probably put it
in the middlemost position, but providing pure white will definitely put
it as the first, brigthest color).`

var VariantsCommand cli.Command = cli.Command{
	Name:                  "variants",
	Aliases:               []string{"v"},
	Usage:                 "generate a scale of variants from a given color",
	ArgsUsage:             "<color>",
	Description:           variantsCommandDescription,
	EnableShellCompletion: true,
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:  "color",
			Value: "",
		},
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "format",
			Value:   "ascii",
			Aliases: []string{"f"},
			Usage:   "output the generated palette in a specific format",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		cmdCfg := &VariantsCommandOptions{
			Color:  c.StringArg("color"),
			Format: c.String("format"),
		}

		return ExecuteVariants(cmdCfg)
	},
}

func ExecuteVariants(cmdCfg *VariantsCommandOptions) error {
	rgb := input.RGBfromString(cmdCfg.Color)

	lab, err := colors.RGBtoLab(rgb)
	if err != nil {
		return fmt.Errorf("could not convert to lab: %w", err)
	}

	lch := colors.LabToLCH(lab)
	vars := variants.GenerateVariants(lch)

	var formatter format.VariantsFormatter
	switch cmdCfg.Format {
	case "ascii":
		formatter = &format.ASCIIVariantsFormatter{}
	case "toml":
		formatter = &format.TOMLVariantsFormatter{}
	default:
		return fmt.Errorf("unknown format specifier")
	}

	fmt.Print(formatter.Format(vars))

	return nil
}
