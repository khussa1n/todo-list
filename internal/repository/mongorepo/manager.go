package mongorepo

import (
	"github.com/khussa1n/todo-list/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	taskCollection *mongo.Collection
}

func New(db *mongo.Database, collections config.Collections) *MongoDB {
	return &MongoDB{
		taskCollection: db.Collection(collections.Task),
	}
}
