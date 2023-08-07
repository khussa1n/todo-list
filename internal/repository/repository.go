package repository

import (
	"context"
	"github.com/khussa1n/todo-list/internal/entity"
)

type TodoList interface {
	CreateTask(ctx context.Context, e *entity.Tasks) error
	UpdateTask(ctx context.Context, e *entity.Tasks) error
	GetAllTasks(ctx context.Context, status string) (*entity.Tasks, error)
	DeleteTask(ctx context.Context, id int64) error
}

type Repository interface {
	TodoList
}
