package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type Repository struct {
	*firestore.Client
}

var fireDB Repository

func (r *Repository) Connect() error {
	ctx := context.Background()
	home, err := os.Getwd()
	if err != nil {
		return err
	}
	opt := option.WithCredentialsFile(home + "/key.json")
	//config := &firestore.Config{DatabaseURL: "https://listtogether-final-default-rtdb.firebaseio.com/"}
	app, err := firestore.NewClient(ctx, "listtogether-final", opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	//client, err := app.Firestore(ctx)
	//if err != nil {
	//	return fmt.Errorf("error initializing database: %v", err)
	//}
	r.Client = app
	return nil
}

func FirebaseDB() *Repository {
	return &fireDB
}

//https://firebase.google.com/docs/firestore/query-data/queries

func (r *Repository) FindFirst(collectionName string, propName string, propValue string, operator string, ctx *gin.Context) (map[string]interface{}, error) {
	iter := r.Client.Collection(collectionName).Where(propName, operator, propValue).Documents(ctx)
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		return doc.Data(), nil
	}
	return nil, nil
}

func (r *Repository) FindAll(collectionName string, propName string, propValue string, operator string, ctx *gin.Context) ([]map[string]interface{}, error) {
	iter := r.Client.Collection(collectionName).Where(propName, operator, propValue).Documents(ctx)
	list := make([]map[string]interface{}, 0)
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		list = append(list, doc.Data())
	}
	return list, nil
}

func (r *Repository) GetById(collectionName string, id string, ctx *gin.Context) (map[string]interface{}, error) {
	dsnap, err := r.Client.Collection(collectionName).Doc(id).Get(ctx)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return nil, nil
			}
		}
		return nil, err
	}

	return dsnap.Data(), nil
}
func (r *Repository) Create(path string, id string, u interface{}, ctx *gin.Context) error {
	if _, err := r.Client.Collection(path).Doc(id).Set(ctx, u); err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(path string, id string, ctx *gin.Context) error {
	_, err := r.Client.Collection(path).Doc(id).Delete(ctx)
	return err
}
