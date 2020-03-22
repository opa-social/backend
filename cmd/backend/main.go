package main

import (
	"os"

	"github.com/opa-social/backend/internal/database"
	"github.com/opa-social/backend/internal/firebase"
	"github.com/opa-social/backend/internal/router"
	"github.com/urfave/cli/v2"
)

func run(ctx *cli.Context) error {
	controller := firebase.New()
	redis := database.New(ctx.String("datastore"), ctx.String("password"))

	router := router.Setup(ctx.String("address"), ctx.Uint("port"), &controller, &redis)
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
			&cli.StringFlag{
				Name:    "datastore",
				Aliases: []string{"d"},
				Usage:   "Redis datastore url.",
				Value:   "redis://0.0.0.0:6379",
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "Redis datastore password (if applicable).",
				Value: "",
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
