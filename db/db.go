package db

import (
	"context"
	"fmt"

	"github.com/cephalization/chihuahua-bot/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDatabaseClient returns a handle onto the running mongo db
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
func NewDatabaseClient() (*mongo.Client, error) {

	// Get credentials from env
	dbName, err := utils.GetEnv("MONGO_INITDB_DATABASE")
	if err != nil {
		return nil, err
	}
	dbUser, err := utils.GetEnv("MONGO_INITDB_ROOT_USERNAME")
	if err != nil {
		return nil, err
	}
	dbPassword, err := utils.GetEnv("MONGO_INITDB_ROOT_PASSWORD")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("mongodb://%s:%s@127.0.0.1:27017/%s", dbUser, dbPassword, dbName)

	// TODO: learn about context
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}
