package service

import (
	"context"
	"github.com/khussa1n/todo-list/internal/entity"
	"github.com/khussa1n/todo-list/internal/entity/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoList interface {
	CreateTask(ctx context.Context, t *dto.TasksDTO) (*entity.Tasks, error)
	UpdateTask(ctx context.Context, t *dto.TasksDTO, id primitive.ObjectID) error
	UpdateTaskStatus(ctx context.Context, id primitive.ObjectID, status string) error
	GetAllTasks(ctx context.Context, status string) ([]entity.Tasks, error)
	DeleteTask(ctx context.Context, id primitive.ObjectID) error
}

type Service interface {
	TodoList
}
