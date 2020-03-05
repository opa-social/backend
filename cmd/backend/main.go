package main

import (
	"os"

	"github.com/opa-social/backend/internal/firebase"
	"github.com/opa-social/backend/internal/router"
	"github.com/urfave/cli/v2"
)

func run(ctx *cli.Context) error {
	controller := firebase.New()

	router := router.Setup(ctx.String("address"), ctx.Uint("port"), &controller)
	router.Serve()

	return nil
}

func main() {
	app := cli.App{
		Name:    "backend",
		Version: "0.1.0",
		Usage:   "Opa! backend web service.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Usage:   "Local address to bind to.",
				Value:   "0.0.0.0",
			},
			&cli.UintFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Local port number to bind to.",
				Value:   9000,
			},
		},
		Action: run,
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Opa!",
				Email: "contact@opa.social",
			},
		},
	}

	app.Run(os.Args)
}
