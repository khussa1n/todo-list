package service

import (
	"context"
	"errors"
	"github.com/khussa1n/todo-list/internal/entity"
	"github.com/khussa1n/todo-list/internal/entity/dto"
	"log"
	"time"
)

func (m *Manager) CreateTask(ctx context.Context, t *dto.TasksDTO) (*entity.Tasks, error) {
	if len(t.Title) > 200 {
		return nil, errors.New("more than 200 char")
	}

	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, t.ActiveAt)
	if err != nil {
		log.Println("activeAt format err: ", err)
		return nil, errors.New("activeAt invalid format")
	}

	weekday := parsedDate.Weekday()

	var title string
	if weekday == time.Saturday || weekday == time.Sunday {
		title = "ВЫХОДНОЙ - "
		title += t.Title
	} else {
		title += t.Title
	}

	task := &entity.Tasks{
		ID:       "",
		Title:    title,
		ActiveAt: t.ActiveAt,
		Status:   "active",
	}

	ID, err := m.Repository.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	task.ID = ID

	return task, nil
}

func (m *Manager) UpdateTask(ctx context.Context, t *dto.TasksDTO, id string) error {
	task, err := m.Repository.GetTaskByID(ctx, id)
	if err != nil {
		return err
	}

	newTask := &entity.Tasks{
		ID:       id,
		Title:    t.Title,
		ActiveAt: t.ActiveAt,
		Status:   task.Status,
	}

	err = m.Repository.UpdateTask(ctx, newTask)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateTaskStatus(ctx context.Context, id string, status string) error {
	err := m.Repository.UpdateTaskStatus(ctx, id, status)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetAllTasks(ctx context.Context, status string) ([]entity.Tasks, error) {
	if status == "" {
		status = "active"
	}
	tasks, err := m.Repository.GetAllTasks(ctx, status)
	if err != nil {
		return nil, err
	}

	if tasks == nil {
		return make([]entity.Tasks, 0), nil
	}

	return tasks, nil
}

func (m *Manager) DeleteTask(ctx context.Context, id string) error {
	err := m.Repository.DeleteTask(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
