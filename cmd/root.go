package cmd

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func Execute() {
	cmd := &cli.Command{
		Name:                  "chromatika",
		EnableShellCompletion: true,
		Version:               "v0.1.0",
		Authors: []any{
			"Daniel Sanso <cs.daniel.sanso@gmail.com>",
		},
		Copyright:   "(c) 2025 Daniel Sanso",
		Usage:       "dominant color extraction cli",
		HideHelp:    false,
		HideVersion: false,
		Commands: []*cli.Command{
			&ExtractCommand,
			&PaletteCommand,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
