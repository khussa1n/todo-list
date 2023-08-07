package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	host     string
	username string
	password string
	port     string
	dbName   string
	db       *mongo.Database
}

func New(DBopts ...Option) (*mongo.Database, error) {
	m := new(Mongo)

	for _, opt := range DBopts {
		opt(m)
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", m.username, m.password, m.host, m.port)

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("mongoDB connect err: %w", err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	var result bson.M
	if err = client.Database(m.dbName).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return nil, fmt.Errorf("mongoDB Send a ping err: %e", err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client.Database(m.dbName), nil
}
