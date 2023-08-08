package mongorepo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const tasksTable = "tasks"

type MongoDB struct {
	collection *mongo.Collection
}

func New(db *mongo.Database, collection string) *MongoDB {
	return &MongoDB{
		collection: db.Collection(collection),
	}
}
