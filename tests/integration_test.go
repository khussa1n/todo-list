package tests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/khussa1n/todo-list/internal/config"
	"github.com/khussa1n/todo-list/internal/handler"
	"github.com/khussa1n/todo-list/internal/repository/mongorepo"
	"github.com/khussa1n/todo-list/internal/service"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"testing"
)

// До запуска теста поднимите docker-compose
func TestIntegration(t *testing.T) {
	cfg, err := config.InitConfig("../config.yaml")
	require.NoError(t, err)

	// Порт mongodb=27020 как в docker-compose
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.DB.Username, cfg.DB.Password, "localhost", "27020")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	testDB := client.Database(cfg.DB.DBName)

	db := mongorepo.New(testDB, cfg.DB.Collections)
	srvs := service.New(db, cfg)
	hndlr := handler.New(srvs)

	var buf bytes.Buffer

	url := fmt.Sprintf("/api/todo-list/tasks/")
	request, err := http.NewRequest(http.MethodGet, url, &buf)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	hndlr.InitRouter().ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}
