package service

import (
	"context"
	"github.com/khussa1n/todo-list/internal/entity"
)

func (m *Manager) CreateTask(ctx context.Context, e *entity.Tasks) error {
	return nil
}

func (m *Manager) UpdateTask(ctx context.Context, e *entity.Tasks) error {
	return nil
}

func (m *Manager) GetAllTasks(ctx context.Context, status string) (*entity.Tasks, error) {
	return nil, nil
}

func (m *Manager) DeleteTask(ctx context.Context, id int64) error {
	return nil
}
