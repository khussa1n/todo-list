package app

import (
	config "github.com/khussa1n/todo-list/internal/config"
	"github.com/khussa1n/todo-list/internal/handler"
	"github.com/khussa1n/todo-list/internal/repository/mongorepo"
	"github.com/khussa1n/todo-list/internal/service"
	"github.com/khussa1n/todo-list/pkg/client/mongodb"
	"github.com/khussa1n/todo-list/pkg/httpserver"
	"log"
	"os"
	"os/signal"
)

func Run(cfg *config.Config) error {
	// Соеденение с базой
	conn, err := mongodb.New(
		mongodb.WithHost(cfg.DB.Host),
		mongodb.WithPort(cfg.DB.Port),
		mongodb.WithDBName(cfg.DB.DBName),
		mongodb.WithUsername(cfg.DB.Username),
		mongodb.WithPassword(cfg.DB.Password),
	)
	if err != nil {
		log.Printf("connection to mongodb err: %s", err.Error())
		return err
	}
	log.Println("connection success")

	// Получение репозитория <Repository interface> и базы mongodb <MongoDB struct>
	db := mongorepo.New(conn, cfg.DB.Collections)
	// Получение сервиса
	srvs := service.New(db, cfg)
	// Получение контроллера
	hndlr := handler.New(srvs)
	// Создание http сервера
	server := httpserver.New(
		hndlr.InitRouter(),
		httpserver.WithPort(cfg.HTTP.Port),
		httpserver.WithReadTimeout(cfg.HTTP.ReadTimeout),
		httpserver.WithWriteTimeout(cfg.HTTP.WriteTimeout),
		httpserver.WithShutdownTimeout(cfg.HTTP.ShutdownTimeout),
	)

	// Запуск http сервера
	server.Start()
	log.Println("server started")

	// Создание канала для сигналов
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Ожидание сигналов от операционной системы и ошибки от http сервера
	select {
	case s := <-interrupt:
		log.Printf("signal received: %s", s.String())
	case err = <-server.Notify():
		log.Printf("server notify: %s", err.Error())
	}

	// Принудительное остановление сервера
	err = server.Shutdown()
	if err != nil {
		log.Printf("server shutdown err: %s", err)
	}

	return nil
}
