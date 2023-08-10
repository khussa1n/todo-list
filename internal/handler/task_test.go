package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/khussa1n/todo-list/internal/custom_error"
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
	id := primitive.NewObjectID()
	table := []struct {
		name            string
		dto             dto.TasksDTO
		dtoJson         string
		expectedSrvc    entity.Tasks
		expectedSrvcErr error
		httpStatus      int
		responseBody    string
	}{
		{
			name:         "ok ВЫХОДНОЙ",
			dto:          dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-05"},
			expectedSrvc: entity.Tasks{ID: id, Title: "ВЫХОДНОЙ - Купить", ActiveAt: "2023-08-05", Status: "active"},
			httpStatus:   http.StatusCreated,
			responseBody: fmt.Sprintf(`{"id":"%s","title":"ВЫХОДНОЙ - Купить","activeAt":"2023-08-05","status":"active"}`, id.Hex()),
		},
		{
			name:         "ok",
			dto:          dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-04"},
			expectedSrvc: entity.Tasks{ID: id, Title: "Купить", ActiveAt: "2023-08-04", Status: "active"},
			httpStatus:   http.StatusCreated,
			responseBody: fmt.Sprintf(`{"id":"%s","title":"Купить","activeAt":"2023-08-04","status":"active"}`, id.Hex()),
		},
		{
			name:            "activeAt invalid format",
			dto:             dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-32"},
			expectedSrvcErr: custom_error.ErrInvalidActiveAtFormat,
			httpStatus:      http.StatusBadRequest,
			responseBody:    `"activeAt invalid format"`,
		},
		{
			name: "more than 200 char",
			dto: dto.TasksDTO{Title: "Купитьasdfsasddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"dddddddddddddddddddddddddddddddddddddddddddddddadsfdd", ActiveAt: "2023-08-04"},
			expectedSrvcErr: custom_error.ErrMessageTooLong,
			httpStatus:      http.StatusBadRequest,
			responseBody:    `"more than 200 char"`,
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockService := mock_service.NewMockService(controller)

			handler := New(mockService)

			recorder := httptest.NewRecorder()

			switch testCase.name {
			case "ok", "ok ВЫХОДНОЙ":
				mockService.EXPECT().CreateTask(gomock.Any(), &testCase.dto).Return(&testCase.expectedSrvc, nil).Times(1)
				break
			case "activeAt invalid format", "more than 200 char":
				mockService.EXPECT().CreateTask(gomock.Any(), &testCase.dto).Return(nil, testCase.expectedSrvcErr).Times(1)
				break
			}

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(testCase.dto)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/todo-list/tasks/")

			request, err := http.NewRequest(http.MethodPost, url, &buf)
			require.NoError(t, err)

			handler.InitRouter().ServeHTTP(recorder, request)

			require.Equal(t, testCase.httpStatus, recorder.Code)
			require.Equal(t, testCase.responseBody, recorder.Body.String())
		})
	}
}

func Test_updateTask(t *testing.T) {
	table := []struct {
		name            string
		dto             dto.TasksDTO
		id              primitive.ObjectID
		idErr           string
		expectedSrvcErr error
		httpStatus      int
		responseBody    string
	}{
		{
			name:         "ok",
			dto:          dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-05"},
			id:           primitive.NewObjectID(),
			httpStatus:   http.StatusNoContent,
			responseBody: "",
		},
		{
			name:            "empty id param",
			dto:             dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-18"},
			id:              primitive.NilObjectID,
			idErr:           "",
			expectedSrvcErr: custom_error.ErrEmptyID,
			httpStatus:      http.StatusNotFound,
			responseBody:    `404 page not found`,
		},
		{
			name:            "activeAt invalid format",
			dto:             dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-52"},
			id:              primitive.NewObjectID(),
			expectedSrvcErr: custom_error.ErrInvalidActiveAtFormat,
			httpStatus:      http.StatusBadRequest,
			responseBody:    `"activeAt invalid format"`,
		},
		{
			name:            "invalid id param",
			dto:             dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-15"},
			id:              primitive.NilObjectID,
			idErr:           "1234",
			expectedSrvcErr: custom_error.ErrInvalidIDParameter,
			httpStatus:      http.StatusBadRequest,
			responseBody:    `"invalid id param"`,
		},
		{
			name: "more than 200 char",
			dto: dto.TasksDTO{Title: "Купитьasdfsasddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"dddddddddddddddddddddddddddddddddddddddddddddddadsfdd", ActiveAt: "2023-08-05"},
			id:              primitive.NewObjectID(),
			expectedSrvcErr: custom_error.ErrMessageTooLong,
			httpStatus:      http.StatusBadRequest,
			responseBody:    `"more than 200 char"`,
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockService := mock_service.NewMockService(controller)

			handler := New(mockService)

			recorder := httptest.NewRecorder()

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(testCase.dto)
			require.NoError(t, err)

			switch testCase.name {
			case "ok":
				mockService.EXPECT().UpdateTask(gomock.Any(), &testCase.dto, testCase.id).Return(nil).Times(1)
				break
			case "activeAt invalid format", "more than 200 char":
				mockService.EXPECT().UpdateTask(gomock.Any(), &testCase.dto, testCase.id).Return(testCase.expectedSrvcErr).Times(1)
				break
			case "empty id param", "invalid id param":
				mockService.EXPECT().UpdateTask(gomock.Any(), &testCase.dto, testCase.idErr).Return(testCase.expectedSrvcErr).Times(0)
				break
			}

			var url string
			if testCase.id == primitive.NilObjectID {
				url = fmt.Sprintf("/api/todo-list/tasks/" + testCase.idErr)
			} else {
				url = fmt.Sprintf("/api/todo-list/tasks/" + testCase.id.Hex())
			}
			request, err := http.NewRequest(http.MethodPut, url, &buf)
			require.NoError(t, err)

			handler.InitRouter().ServeHTTP(recorder, request)

			require.Equal(t, testCase.httpStatus, recorder.Code)
			require.Equal(t, testCase.responseBody, recorder.Body.String())
		})
	}
}

func Test_updateTaskStatus(t *testing.T) {
	table := []struct {
		name            string
		id              primitive.ObjectID
		idErr           string
		status          string
		expectedSrvcErr error
		httpStatus      int
		responseBody    string
	}{
		{
			name:       "ok",
			status:     "done",
			id:         primitive.NewObjectID(),
			httpStatus: http.StatusNoContent,
		},
		{
			name:            "empty id param",
			id:              primitive.NilObjectID,
			idErr:           "",
			status:          "done",
			expectedSrvcErr: custom_error.ErrEmptyID,
			httpStatus:      http.StatusBadRequest,
			responseBody:    `"empty id param"`,
		},
		{
			name:            "invalid id param",
			id:              primitive.NilObjectID,
			idErr:           "1234",
			status:          "done",
			expectedSrvcErr: custom_error.ErrInvalidIDParameter,
			httpStatus:      http.StatusBadRequest,
			responseBody:    `"invalid id param"`,
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockService := mock_service.NewMockService(controller)

			handler := New(mockService)

			recorder := httptest.NewRecorder()

			switch testCase.name {
			case "ok":
				mockService.EXPECT().UpdateTaskStatus(gomock.Any(), testCase.id, testCase.status).Return(nil).Times(1)
				break
			case "empty id param", "invalid id param":
				mockService.EXPECT().UpdateTaskStatus(gomock.Any(), testCase.id, testCase.status).Return(nil).Times(0)
				break
			}

			var url string
			if testCase.id == primitive.NilObjectID {
				url = fmt.Sprintf("/api/todo-list/tasks/" + testCase.idErr + "/" + testCase.status)
			} else {
				url = fmt.Sprintf("/api/todo-list/tasks/" + testCase.id.Hex() + "/" + testCase.status)
			}
			request, err := http.NewRequest(http.MethodPut, url, nil)
			require.NoError(t, err)

			handler.InitRouter().ServeHTTP(recorder, request)

			require.Equal(t, testCase.httpStatus, recorder.Code)
			require.Equal(t, testCase.responseBody, recorder.Body.String())
		})
	}
}

func Test_deleteTask(t *testing.T) {
	table := []struct {
		name         string
		id           primitive.ObjectID
		idErr        string
		responseBody string
		httpStatus   int
	}{
		{
			name:       "ok",
			id:         primitive.NewObjectID(),
			httpStatus: http.StatusNoContent,
		},
		{
			name:         "empty id param",
			id:           primitive.NilObjectID,
			idErr:        "",
			httpStatus:   http.StatusNotFound,
			responseBody: `404 page not found`,
		},
		{
			name:         "invalid id param",
			id:           primitive.NilObjectID,
			idErr:        "1234",
			httpStatus:   http.StatusBadRequest,
			responseBody: `"invalid id param"`,
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockService := mock_service.NewMockService(controller)

			handler := New(mockService)

			recorder := httptest.NewRecorder()

			switch testCase.name {
			case "ok":
				mockService.EXPECT().DeleteTask(gomock.Any(), testCase.id).Return(nil).Times(1)
				break
			case "empty id param", "invalid id param":
				mockService.EXPECT().DeleteTask(gomock.Any(), testCase.id).Return(nil).Times(0)
				break
			}

			var url string
			if testCase.id == primitive.NilObjectID {
				url = fmt.Sprintf("/api/todo-list/tasks/" + testCase.idErr)
			} else {
				url = fmt.Sprintf("/api/todo-list/tasks/" + testCase.id.Hex())
			}

			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			handler.InitRouter().ServeHTTP(recorder, request)

			require.Equal(t, testCase.httpStatus, recorder.Code)
			require.Equal(t, testCase.responseBody, recorder.Body.String())
		})
	}
}

func Test_getAllTasks(t *testing.T) {
	table := []struct {
		name            string
		status          string
		expectedService []entity.Tasks
		httpStatus      int
	}{
		{
			name:            "ok",
			expectedService: []entity.Tasks{{ID: primitive.NewObjectID(), Title: "Купить", ActiveAt: "2023-08-05", Status: "active"}},
			httpStatus:      http.StatusOK,
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockService := mock_service.NewMockService(controller)

			handler := New(mockService)

			recorder := httptest.NewRecorder()

			mockService.EXPECT().GetAllTasks(gomock.Any(), testCase.status).Return(testCase.expectedService, nil).Times(1)

			url := fmt.Sprintf("/api/todo-list/tasks/" + testCase.status)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			handler.InitRouter().ServeHTTP(recorder, request)

			require.Equal(t, testCase.httpStatus, recorder.Code)
		})
	}
}
