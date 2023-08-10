package tests

import (
	"github.com/khussa1n/todo-list/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	task = entity.Tasks{
		ID:       primitive.NewObjectID(),
		Title:    "Title",
		ActiveAt: "2023-08-04",
		Status:   "active",
	}
)
