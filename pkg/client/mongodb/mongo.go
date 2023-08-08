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
}

func New(cfgOpts ...Option) (*mongo.Database, error) {
	m := new(Mongo)

	for _, opt := range cfgOpts {
		opt(m)
	}

	url := fmt.Sprintf("mongodb://%s:%s@%s:%s", m.username, m.password, m.host, m.port)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)

	// Создание нового клиент и подключение к db
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("mongoDB connect err: %w", err)
	}

	// ping-запрос для подтверждения успешного подключения
	var result bson.M
	if err = client.Database(m.dbName).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return nil, fmt.Errorf("mongoDB Send a ping err: %e", err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client.Database(m.dbName), nil
}
