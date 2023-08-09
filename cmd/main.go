package main

import (
	"github.com/khussa1n/todo-list/internal/app"
	"github.com/khussa1n/todo-list/internal/config"
)

// @title           Todo List
// @version         0.0.1
// @description     API for Todo application.

// @contact.name   Khussain
// @contact.email  khussain.qudaibergenov@gmail.com

// @host      localhost:8080
// @BasePath  /api/todo-list
func main() {
	// Инициализация кофигурации
	cfg, err := config.InitConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	// Запуск
	err = app.Run(cfg)
	if err != nil {
		panic(err)
	}
}
