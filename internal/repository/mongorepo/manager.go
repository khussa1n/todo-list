package mongorepo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const tasksTable = "tasks"

type MongoDB struct {
	DB *mongo.Database
}

func New(db *mongo.Database) *MongoDB {
	return &MongoDB{
		DB: db,
	}
}
