package main

import (
	"context"
	"github.com/alhamsya/voltron/cmd/rest"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	cliApp := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run protocol service",
				Subcommands: []*cli.Command{
					{
						Name:  "rest",
						Usage: "run rest API",
						Action: func(ctx *cli.Context) error {
							return rest.RunApp(ctx.Context)
						},
					},
				},
			},
		},
		AllowExtFlags: true,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name: "cfg-credential",
			},
			&cli.IntFlag{
				Name: "cfg-static",
			},
		},
	}

	if err := cliApp.RunContext(context.Background(), os.Args); err != nil {
		panic(err.Error())
	}
}
