package repository

import (
	"context"
	"firebase.google.com/go/db"
	"fmt"
	"os"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

type Repository struct {
	*db.Client
}

var fireDB Repository

func (db *Repository) Connect() error {
	home, err := os.Getwd()
	if err != nil {
		return err
	}
	ctx := context.Background()
	opt := option.WithCredentialsFile(home + "/key.json")
	config := &firebase.Config{DatabaseURL: "https://listtogether-final-default-rtdb.firebaseio.com/"}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	client, err := app.Database(ctx)
	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}
	db.Client = client
	return nil
}

func FirebaseDB() *Repository {
	return &fireDB
}
