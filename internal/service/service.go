package service

import (
	"context"
	"github.com/khussa1n/todo-list/internal/entity"
	"github.com/khussa1n/todo-list/internal/entity/dto"
)

type TodoList interface {
	CreateTask(ctx context.Context, t *dto.TasksDTO) (*entity.Tasks, error)
	UpdateTask(ctx context.Context, t *dto.TasksDTO, id string) error
	UpdateTaskStatus(ctx context.Context, id string, status string) error
	GetAllTasks(ctx context.Context, status string) ([]entity.Tasks, error)
	DeleteTask(ctx context.Context, id string) error
}

type Service interface {
	TodoList
}
