package main

import (
	"context"
	"fmt"
	"os"

	"github.com/opa-social/backend/internal/firebase"
	"github.com/urfave/cli/v2"
)

func run(ctx *cli.Context) error {
	controller := firebase.New()

	token, err := controller.Client.CustomToken(context.Background(), ctx.String("uid"))
	if err != nil {
		cli.Exit(err, 1)
	}

	fmt.Println(token)
	return nil
}

func main() {
	app := cli.App{
		Name:    "get-token",
		Version: "0.1.0",
		Usage:   "Tool to get JWT for user with specific UID.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "uid",
				Usage:    "User's uid",
				Required: true,
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
