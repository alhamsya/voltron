package main

import (
	"context"
	"os"

	"github.com/alhamsya/voltron/cmd/consumer"
	"github.com/alhamsya/voltron/cmd/rest"
	"github.com/urfave/cli/v2"
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
					{
						Name:  "consumer",
						Usage: "run consumer",
						Action: func(ctx *cli.Context) error {
							return consumer.RunApp(ctx.Context)
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
