package mongorepo

import (
	"context"
	"github.com/khussa1n/todo-list/internal/entity"
)

func (m *MongoDB) CreateTask(ctx context.Context, e *entity.Tasks) error {
	return nil
}

func (m *MongoDB) UpdateTask(ctx context.Context, e *entity.Tasks) error {
	return nil
}

func (m *MongoDB) GetAllTasks(ctx context.Context, status string) (*entity.Tasks, error) {
	return nil, nil
}

func (m *MongoDB) DeleteTask(ctx context.Context, id int64) error {
	return nil
}
