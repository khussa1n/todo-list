package main

import (
	"fmt"
	"github.com/khussa1n/todo-list/internal/app"
	"github.com/khussa1n/todo-list/internal/config"
)

func main() {
	cfg, err := config.InitConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%#v", cfg))

	err = app.Run(cfg)
	if err != nil {
		panic(err)
	}
}
