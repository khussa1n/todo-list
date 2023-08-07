package service

import (
	"github.com/khussa1n/todo-list/internal/config"
	"github.com/khussa1n/todo-list/internal/repository"
)

type Manager struct {
	Repository repository.Repository
	Config     *config.Config
}

func New(repository repository.Repository, config *config.Config) *Manager {
	return &Manager{
		Repository: repository,
		Config:     config,
	}
}
