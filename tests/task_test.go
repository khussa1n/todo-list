package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/khussa1n/todo-list/internal/entity"
	"github.com/khussa1n/todo-list/internal/entity/dto"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestCreateTask() {
	dto := dto.TasksDTO{Title: "test", ActiveAt: "2023-08-04"}

	router := s.handler.InitRouter()

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(dto)
	s.NoError(err)

	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/api/todo-list/tasks/")
	request, err := http.NewRequest(http.MethodPost, url, &buf)
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Referer", "https://workshop.zhashkevych.com/")
	s.NoError(err)

	router.ServeHTTP(recorder, request)

	r := s.Require()

	if recorder.Code == 400 {
		r.Equal("\"a task with the same title already exists\"", recorder.Body.String())
	} else {
		r.Equal(http.StatusCreated, recorder.Code)
		var tasks entity.Tasks
		err = s.db.Collection(cfg.Test.DB.Collections.Task).FindOne(context.Background(), bson.M{"title": dto.Title}).Decode(&tasks)
		s.NoError(err)

		r.Equal(dto.Title, tasks.Title)
		r.Equal(dto.ActiveAt, tasks.ActiveAt)
		r.Equal("active", tasks.Status)
	}
}

func (s *APITestSuite) TestGetAllTasks() {
	var buf bytes.Buffer

	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/api/todo-list/tasks/")
	request, err := http.NewRequest(http.MethodGet, url, &buf)
	s.NoError(err)

	s.handler.InitRouter().ServeHTTP(recorder, request)

	r := s.Require()

	r.Equal(http.StatusOK, recorder.Code)
}
