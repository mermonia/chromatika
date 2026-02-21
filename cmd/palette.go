package cmd

import (
	"github.com/urfave/cli/v3"
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
	Commands: []*cli.Command{
		&PaletteImageCommand,
	},
}
