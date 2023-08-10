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
	dtoTest := dto.TasksDTO{Title: "test", ActiveAt: "2023-08-04"}

	router := s.handler.InitRouter()

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(dtoTest)
	s.NoError(err)

	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/api/todo-list/tasks/")
	request, err := http.NewRequest(http.MethodPost, url, &buf)
	s.NoError(err)

	router.ServeHTTP(recorder, request)

	r := s.Require()

	if recorder.Code == 400 {
		r.Equal("\"a task with the same title already exists\"", recorder.Body.String())
	} else {
		r.Equal(http.StatusCreated, recorder.Code)
		var task2 entity.Tasks
		err = s.db.Collection(cfg.Test.DB.Collections.Task).FindOne(context.Background(), bson.M{"title": dtoTest.Title}).Decode(&task2)
		s.NoError(err)

		r.Equal(dtoTest.Title, task2.Title)
		r.Equal(dtoTest.ActiveAt, task2.ActiveAt)
		r.Equal("active", task2.Status)
	}
}

func (s *APITestSuite) TestUpdateTask() {
	_, err = s.db.Collection(cfg.DB.Collections.Task).InsertOne(context.Background(), task)
	s.NoError(err)

	dtoTest := dto.TasksDTO{Title: "test2_update", ActiveAt: "2023-08-04"}

	router := s.handler.InitRouter()

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(dtoTest)
	s.NoError(err)

	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/api/todo-list/tasks/" + task.ID.Hex())
	request, err := http.NewRequest(http.MethodPut, url, &buf)
	s.NoError(err)

	router.ServeHTTP(recorder, request)

	r := s.Require()

	r.Equal(http.StatusNoContent, recorder.Code)
	var task2 entity.Tasks
	err = s.db.Collection(cfg.Test.DB.Collections.Task).FindOne(context.Background(), bson.M{"_id": task.ID}).Decode(&task2)
	s.NoError(err)

	r.Equal(dtoTest.Title, task2.Title)
	r.Equal(dtoTest.ActiveAt, task2.ActiveAt)
}

func (s *APITestSuite) TestUpdateTaskStatus() {
	router := s.handler.InitRouter()

	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/api/todo-list/tasks/" + task.ID.Hex() + "/done")
	request, err := http.NewRequest(http.MethodPut, url, nil)
	s.NoError(err)

	router.ServeHTTP(recorder, request)

	r := s.Require()

	r.Equal(http.StatusNoContent, recorder.Code)
	var task2 entity.Tasks
	err = s.db.Collection(cfg.Test.DB.Collections.Task).FindOne(context.Background(), bson.M{"_id": task.ID}).Decode(&task2)
	s.NoError(err)

	r.Equal("done", task2.Status)
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

func (s *APITestSuite) TestDeleteTask() {
	var buf bytes.Buffer

	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/api/todo-list/tasks/" + task.ID.Hex())
	request, err := http.NewRequest(http.MethodDelete, url, &buf)
	s.NoError(err)

	s.handler.InitRouter().ServeHTTP(recorder, request)

	r := s.Require()

	r.Equal(http.StatusNoContent, recorder.Code)
}
