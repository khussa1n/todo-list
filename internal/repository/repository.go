package repository

import (
	"context"
	"github.com/khussa1n/todo-list/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoList interface {
	CreateTask(ctx context.Context, e *entity.Tasks) (*entity.Tasks, error)
	UpdateTask(ctx context.Context, e *entity.Tasks, id primitive.ObjectID) error
	UpdateTaskStatus(ctx context.Context, id primitive.ObjectID, status string) error
	GetAllTasks(ctx context.Context, status string) ([]entity.Tasks, error)
	GetTaskByID(ctx context.Context, id primitive.ObjectID) (*entity.Tasks, error)
	DeleteTask(ctx context.Context, id primitive.ObjectID) error
}

type Repository interface {
	TodoList
}
