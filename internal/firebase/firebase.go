package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
)

// Controller contains app and client instances for Firebase related requests.
type Controller struct {
	// Instance is the instance of the Firebase app.
	Instance *firebase.App
	// Client is the authentication instance used by the controller.
	Client *auth.Client
	// Database is the database instance used by the controller.
	Database *db.Client
}

// New creates a new Firebase controller. Requires that GOOGLE_APPLICATION_CREDENTIALS
// environment variable be set.
func New() Controller {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("Could not initialize Firebase SDK: \"%s\"", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Could not authenticate with Firebase auth:  \"%s\"", err)
	}

	db, err := app.Database(context.Background())
	if err != nil {
		log.Fatalf("Could not authenticate with Firebase database: \"%s\"", err)
	}

	return Controller{
		Instance: app,
		Client:   client,
		Database: db,
	}
}
