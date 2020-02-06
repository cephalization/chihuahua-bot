package db

import (
	"context"
	"fmt"

	"github.com/cephalization/chihuahua-bot/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDatabase returns a handle onto the running mongo db
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
// https://godoc.org/go.mongodb.org/mongo-driver/mongo
func NewDatabase() (*mongo.Database, error) {

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

	uri := fmt.Sprintf("mongodb://%s:%s@database:27017", dbUser, dbPassword)

	// TODO: learn about context
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	DB := client.Database(dbName)

	return DB, nil
}
