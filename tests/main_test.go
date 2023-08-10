package tests

import (
	"context"
	"github.com/khussa1n/todo-list/internal/config"
	"github.com/khussa1n/todo-list/internal/handler"
	"github.com/khussa1n/todo-list/internal/repository/mongorepo"
	"github.com/khussa1n/todo-list/internal/service"
	"github.com/khussa1n/todo-list/pkg/client/mongodb"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
)

var (
	cfg *config.Config
	err error
)

func init() {
	cfg, err = config.InitConfig("../config.yaml")
	if err != nil {
		panic(err)
	}
}

type APITestSuite struct {
	suite.Suite

	db      *mongo.Database
	handler *handler.Handler
	service *service.Manager
	repos   *mongorepo.MongoDB
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	if client, err := mongodb.New(
		mongodb.WithHost(cfg.Test.DB.Host),
		mongodb.WithPort(cfg.Test.DB.Port),
		mongodb.WithDBName(cfg.Test.DB.DBName),
	); err != nil {
		s.FailNow("Failed to connect to mongo", err)
	} else {
		s.db = client
	}

	s.initDeps()

	if err = s.populateDB(); err != nil {
		s.FailNow("Failed to populate DB", err)
	}

}

func (s *APITestSuite) TearDownSuite() {
	s.db.Client().Disconnect(context.Background())
}

func (s *APITestSuite) initDeps() {
	mdb := mongorepo.New(s.db, cfg.DB.Collections)
	srvs := service.New(mdb, cfg)
	hndlr := handler.New(srvs)

	s.repos = mdb
	s.service = srvs
	s.handler = hndlr
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func (s *APITestSuite) populateDB() error {
	_, err = s.db.Collection(cfg.DB.Collections.Task).InsertOne(context.Background(), task)
	if err != nil {
		return err
	}
	return nil
}
