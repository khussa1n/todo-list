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
	conn, err := mongodb.New(
		mongodb.WithHost(cfg.DB.Host),
		mongodb.WithPort(cfg.DB.Port),
		mongodb.WithDBName(cfg.DB.DBName),
		mongodb.WithUsername(cfg.DB.Username),
		mongodb.WithPassword(cfg.DB.Password),
	)
	if err != nil {
		log.Printf("connection to DB err: %s", err.Error())
		return err
	}
	log.Println("connection success")

	db := mongorepo.New(conn)

	srvs := service.New(db, cfg)
	hndlr := handler.New(srvs)
	server := httpserver.New(
		hndlr.InitRouter(),
		httpserver.WithPort(cfg.HTTP.Port),
		httpserver.WithReadTimeout(cfg.HTTP.ReadTimeout),
		httpserver.WithWriteTimeout(cfg.HTTP.WriteTimeout),
		httpserver.WithShutdownTimeout(cfg.HTTP.ShutdownTimeout),
	)

	log.Println("server started")
	server.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case s := <-interrupt:
		log.Printf("signal received: %s", s.String())
	case err = <-server.Notify():
		log.Printf("server notify: %s", err.Error())
	}

	err = server.Shutdown()
	if err != nil {
		log.Printf("server shutdown err: %s", err)
	}

	return nil
}
