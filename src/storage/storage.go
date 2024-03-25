package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	Connection     *mongo.Client
	UserCollection *mongo.Collection
	BookCollection *mongo.Collection
)

func GetConnection() (*mongo.Client, error) {
	if Connection == nil {
		clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			return nil, err
		}

		err = client.Ping(context.Background(), nil)
		if err != nil {
			return nil, err
		}

		fmt.Println("Connected to MongoDB!")
		Connection = client
		BookCollection = Connection.Database("test").Collection("book")
		_, err = BookCollection.Indexes().CreateOne(context.Background(),
			mongo.IndexModel{
				Keys:    bson.D{{"id", 1}},
				Options: options.Index().SetUnique(true),
			})
		if err != nil {
			log.Fatalf("error during setting index model %v \n", err)
		}
		UserCollection = Connection.Database("test").Collection("users")

		_, err = UserCollection.Indexes().CreateMany(context.Background(),
			[]mongo.IndexModel{
				{
					Keys:    bson.D{{"login", 1}},
					Options: options.Index().SetUnique(true),
				}, {
					Keys:    bson.D{{"id", 1}},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys:    bson.D{{"email", 1}},
					Options: options.Index().SetUnique(true),
				},
			})
	}
	return Connection, nil
}
