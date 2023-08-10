package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/khussa1n/todo-list/internal/entity"
	"github.com/khussa1n/todo-list/internal/entity/dto"
	mock_service "github.com/khussa1n/todo-list/internal/service/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_createTask(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mock_service.NewMockService(controller)

	handler := New(mockService)

	recorder := httptest.NewRecorder()

	table := []struct {
		dto          dto.TasksDTO
		expectedSrvc entity.Tasks
		httpStatus   int
	}{
		{
			dto:          dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-05"},
			expectedSrvc: entity.Tasks{ID: primitive.NewObjectID(), Title: "ВЫХОДНОЙ - Купить", ActiveAt: "2023-08-05", Status: "active"},
			httpStatus:   http.StatusCreated,
		},
		{
			dto:          dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-04"},
			expectedSrvc: entity.Tasks{ID: primitive.NewObjectID(), Title: "Купить", ActiveAt: "2023-08-04", Status: "active"},
			httpStatus:   http.StatusCreated,
		},
	}

	for _, testCase := range table {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(testCase.dto)
		require.NoError(t, err)

		mockService.EXPECT().CreateTask(gomock.Any(), &testCase.dto).Return(&testCase.expectedSrvc, nil).Times(1)

		url := fmt.Sprintf("/api/todo-list/tasks/")
		request, err := http.NewRequest(http.MethodPost, url, &buf)
		require.NoError(t, err)

		handler.InitRouter().ServeHTTP(recorder, request)

		require.Equal(t, testCase.httpStatus, recorder.Code)
	}
}

func Test_updateTask(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mock_service.NewMockService(controller)

	handler := New(mockService)

	recorder := httptest.NewRecorder()

	table := []struct {
		dto        dto.TasksDTO
		id         primitive.ObjectID
		httpStatus int
	}{
		{
			dto:        dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-05"},
			id:         primitive.NewObjectID(),
			httpStatus: http.StatusNoContent,
		},
	}

	for _, testCase := range table {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(testCase.dto)
		require.NoError(t, err)

		mockService.EXPECT().UpdateTask(gomock.Any(), &testCase.dto, testCase.id).Return(nil).Times(1)

		url := fmt.Sprintf("/api/todo-list/tasks/" + testCase.id.Hex())
		request, err := http.NewRequest(http.MethodPut, url, &buf)
		require.NoError(t, err)

		handler.InitRouter().ServeHTTP(recorder, request)

		require.Equal(t, testCase.httpStatus, recorder.Code)
	}
}

func Test_updateTaskStatus(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mock_service.NewMockService(controller)

	handler := New(mockService)

	recorder := httptest.NewRecorder()

	table := []struct {
		status     string
		id         primitive.ObjectID
		httpStatus int
	}{
		{
			status:     "done",
			id:         primitive.NewObjectID(),
			httpStatus: http.StatusNoContent,
		},
	}

	for _, testCase := range table {
		var buf bytes.Buffer

		mockService.EXPECT().UpdateTaskStatus(gomock.Any(), testCase.id, testCase.status).Return(nil).Times(1)

		url := fmt.Sprintf("/api/todo-list/tasks/" + testCase.id.Hex() + "/" + testCase.status)
		request, err := http.NewRequest(http.MethodPut, url, &buf)
		require.NoError(t, err)

		handler.InitRouter().ServeHTTP(recorder, request)

		require.Equal(t, testCase.httpStatus, recorder.Code)
	}
}

func Test_deleteTask(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mock_service.NewMockService(controller)

	handler := New(mockService)

	recorder := httptest.NewRecorder()

	table := []struct {
		id         primitive.ObjectID
		httpStatus int
	}{
		{
			id:         primitive.NewObjectID(),
			httpStatus: http.StatusNoContent,
		},
	}

	for _, testCase := range table {
		var buf bytes.Buffer

		mockService.EXPECT().DeleteTask(gomock.Any(), testCase.id).Return(nil).Times(1)

		url := fmt.Sprintf("/api/todo-list/tasks/" + testCase.id.Hex())
		request, err := http.NewRequest(http.MethodDelete, url, &buf)
		require.NoError(t, err)

		handler.InitRouter().ServeHTTP(recorder, request)

		require.Equal(t, testCase.httpStatus, recorder.Code)
	}
}

func Test_getAllTasks(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mock_service.NewMockService(controller)

	handler := New(mockService)

	recorder := httptest.NewRecorder()

	table := []struct {
		status          string
		expectedService []entity.Tasks
		httpStatus      int
	}{
		{
			expectedService: []entity.Tasks{{ID: primitive.NewObjectID(), Title: "Купить", ActiveAt: "2023-08-05", Status: "active"}},
			httpStatus:      http.StatusOK,
		},
	}

	for _, testCase := range table {
		var buf bytes.Buffer

		mockService.EXPECT().GetAllTasks(gomock.Any(), testCase.status).Return(testCase.expectedService, nil).Times(1)

		url := fmt.Sprintf("/api/todo-list/tasks/" + testCase.status)
		request, err := http.NewRequest(http.MethodGet, url, &buf)
		require.NoError(t, err)

		handler.InitRouter().ServeHTTP(recorder, request)

		require.Equal(t, testCase.httpStatus, recorder.Code)
	}
}
