package repository

import (
	"context"
	"github.com/khussa1n/todo-list/internal/entity"
)

type TodoList interface {
	CreateTask(ctx context.Context, e *entity.Tasks) (string, error)
	UpdateTask(ctx context.Context, e *entity.Tasks) error
	UpdateTaskStatus(ctx context.Context, id string, status string) error
	GetAllTasks(ctx context.Context, status string) ([]entity.Tasks, error)
	GetTaskByID(ctx context.Context, id string) (*entity.Tasks, error)
	DeleteTask(ctx context.Context, id string) error
}

type Repository interface {
	TodoList
}
