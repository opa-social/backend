package main

import (
	"context"
	"fmt"
	"os"

	"firebase.google.com/go/auth"
	"github.com/opa-social/backend/internal/firebase"
	"github.com/urfave/cli/v2"
)

var userMode = map[string]bool{
	"guest":   true,
	"regular": true,
	"manager": true,
}

func validateMode(mode string) error {
	if _, ok := userMode[mode]; !ok {
		return fmt.Errorf("User type %s not valid", mode)
	}

	return nil
}

func run(ctx *cli.Context) error {
	controller := firebase.New()

	// Ensure that the specified user type is valid.
	err := validateMode(ctx.String("user-type"))
	if err != nil {
		cli.Exit(err, 1)
	}

	userParams := (&auth.UserToCreate{}).
		DisplayName(ctx.String("name")).
		Email(ctx.String("email")).
		Password(ctx.String("password"))

	// Add the user.
	user, err := controller.Client.CreateUser(context.Background(), userParams)
	if err != nil {
		cli.Exit(err, 1)
	}

	// Set the user type using custom claims.
	claims := map[string]interface{}{ctx.String("user-type"): true}
	err = controller.Client.SetCustomUserClaims(context.Background(), user.UID, claims)
	if err != nil {
		cli.Exit(err, 1)
	}

	// Create realtime database entry.
	err = controller.Database.NewRef(fmt.Sprintf("/users/%s", user.UID)).
		Set(context.Background(), map[string]interface{}{"name": ctx.String("name")})
	if err != nil {
		cli.Exit(err, 1)
	}

	fmt.Printf("User %s successfully created\n", ctx.String("name"))
	return nil
}

func main() {
	app := &cli.App{
		Name:    "user-loader",
		Version: "0.1.0",
		Usage:   "Tool to create a user in the Firebase authentication store.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "User's full name.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "email",
				Usage:    "User's email address.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Usage:    "User's password.",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "user-type",
				Usage: "The type of user. One of guest, regular, or manager.",
				Value: "regular",
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
