package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/iamhectorsosa/porygon/internal/mapper"
	"github.com/iamhectorsosa/porygon/internal/render"
	"github.com/urfave/cli/v3"
)

func main() {
	if err := app().Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func app() *cli.Command {
	return &cli.Command{
		Name:   "porygon",
		Usage:  "map a codebase",
		Writer: os.Stdout,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Usage:   "directory to map",
				Value:   ".",
			},
			&cli.StringSliceFlag{
				Name:    "include",
				Aliases: []string{"i"},
				Usage:   "file types to include",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg := mapper.Options{
				Dir:     cmd.String("dir"),
				Include: cmd.StringSlice("include"),
			}
			return run(cfg, cmd.Root().Writer)
		},
	}
}

func run(cfg mapper.Options, stdout io.Writer) error {
	res, err := mapper.Map(cfg)
	if err != nil {
		return err
	}
	return render.Write(stdout, res)
}
