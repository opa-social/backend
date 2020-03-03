package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

// Controller contains app and client instances for Firebase related requests.
type Controller struct {
	// Instance is the instance of the Firebase app.
	Instance *firebase.App
	// Client is the authentication instances used by the controller.
	Client *auth.Client
}

// New creates a new Firebase controller. Requires that GOOGLE_APPLICATION_CREDENTIALS
// environment variable be set.
func New() Controller {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatal("Could not initialize Firebase SDK")
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal("Could not authenticate with Firebase API")
	}

	return Controller{
		Instance: app,
		Client:   client,
	}
}
