package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/khussa1n/todo-list/internal/config"
	"github.com/khussa1n/todo-list/internal/entity"
	"github.com/khussa1n/todo-list/internal/entity/dto"
	mock_repository "github.com/khussa1n/todo-list/internal/repository/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func Test_CreateTask(t *testing.T) {
	cfg, err := config.InitConfig("../../config.yaml")
	require.NoError(t, err)

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_repository.NewMockRepository(controller)

	ctx := context.Background()

	id := primitive.NewObjectID()

	table := []struct {
		dto          dto.TasksDTO
		taskRepo     entity.Tasks
		expectedRepo entity.Tasks
		expectedSrvc entity.Tasks
	}{
		{
			dto:          dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-05"},
			taskRepo:     entity.Tasks{Title: "ВЫХОДНОЙ - Купить", ActiveAt: "2023-08-05", Status: "active"},
			expectedRepo: entity.Tasks{ID: id, Title: "ВЫХОДНОЙ - Купить", ActiveAt: "2023-08-05", Status: "active"},
			expectedSrvc: entity.Tasks{ID: id, Title: "ВЫХОДНОЙ - Купить", ActiveAt: "2023-08-05", Status: "active"},
		},
		{
			dto:          dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-04"},
			taskRepo:     entity.Tasks{Title: "Купить", ActiveAt: "2023-08-04", Status: "active"},
			expectedRepo: entity.Tasks{ID: id, Title: "Купить", ActiveAt: "2023-08-04", Status: "active"},
			expectedSrvc: entity.Tasks{ID: id, Title: "Купить", ActiveAt: "2023-08-04", Status: "active"},
		},
	}

	errorTable := []struct {
		dto             dto.TasksDTO
		expectedSrvcErr string
	}{
		{
			dto:             dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-32"},
			expectedSrvcErr: "activeAt invalid format",
		},
		{
			dto: dto.TasksDTO{Title: "Купитьasdfsasddddddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd", ActiveAt: "2023-08-05"},
			expectedSrvcErr: "more than 200 char",
		},
	}

	for _, testCase := range table {
		mockRepo.EXPECT().CreateTask(ctx, &testCase.taskRepo).Return(&testCase.expectedRepo, nil).Times(1)

		service := New(mockRepo, cfg)

		result, err := service.CreateTask(ctx, &testCase.dto)
		require.NoError(t, err)
		require.NotEmpty(t, result)

		if *result != testCase.expectedSrvc {
			t.Errorf("Incorrect result. Expect %+v, got %+v", testCase.expectedSrvc, result)
		}
	}

	for _, testCase := range errorTable {
		service := New(mockRepo, cfg)

		result, err := service.CreateTask(ctx, &testCase.dto)
		require.Error(t, err)
		require.Empty(t, result)

		if err.Error() != testCase.expectedSrvcErr {
			t.Errorf("Incorrect result. Expect %+v, got %+v", testCase.expectedSrvcErr, err)
		}
	}
}

func Test_UpdateTask(t *testing.T) {
	cfg, err := config.InitConfig("../../config.yaml")
	require.NoError(t, err)

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_repository.NewMockRepository(controller)

	ctx := context.Background()

	id := primitive.NewObjectID()

	table := []struct {
		dto           dto.TasksDTO
		taskRepo      primitive.ObjectID
		expectedRepo  entity.Tasks
		expectedRepo2 entity.Tasks
	}{
		{
			dto:           dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-05"},
			taskRepo:      id,
			expectedRepo:  entity.Tasks{ID: id, Title: "Купить", ActiveAt: "2023-08-05", Status: "active"},
			expectedRepo2: entity.Tasks{Title: "Купить", ActiveAt: "2023-08-05", Status: "active"},
		},
	}

	for _, testCase := range table {
		mockRepo.EXPECT().GetTaskByID(ctx, testCase.taskRepo).Return(&testCase.expectedRepo, nil).Times(1)
		mockRepo.EXPECT().UpdateTask(ctx, &testCase.expectedRepo2, testCase.taskRepo).Return(nil).Times(1)

		service := New(mockRepo, cfg)

		err = service.UpdateTask(ctx, &testCase.dto, testCase.expectedRepo.ID)
		require.NoError(t, err)
	}
}

func Test_UpdateTaskStatus(t *testing.T) {
	cfg, err := config.InitConfig("../../config.yaml")
	require.NoError(t, err)

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_repository.NewMockRepository(controller)

	ctx := context.Background()

	table := []struct {
		id     primitive.ObjectID
		status string
	}{
		{
			id:     primitive.NewObjectID(),
			status: "done",
		},
	}

	for _, testCase := range table {
		mockRepo.EXPECT().UpdateTaskStatus(ctx, testCase.id, testCase.status).Return(nil).Times(1)

		service := New(mockRepo, cfg)

		err = service.UpdateTaskStatus(ctx, testCase.id, testCase.status)
		require.NoError(t, err)
	}
}

func Test_GetAllTasks(t *testing.T) {
	cfg, err := config.InitConfig("../../config.yaml")
	require.NoError(t, err)

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_repository.NewMockRepository(controller)

	ctx := context.Background()

	table := []struct {
		status       string
		expectedRepo []entity.Tasks
	}{
		{
			status:       "active",
			expectedRepo: []entity.Tasks{{ID: primitive.NewObjectID(), Title: "Купить", ActiveAt: "2023-08-05", Status: "active"}},
		},
	}

	for _, testCase := range table {
		mockRepo.EXPECT().GetAllTasks(ctx, testCase.status).Return(testCase.expectedRepo, nil).Times(1)

		service := New(mockRepo, cfg)

		result, err := service.GetAllTasks(ctx, testCase.status)
		require.NoError(t, err)
		require.NotEmpty(t, result)

		if result[0] != testCase.expectedRepo[0] {
			t.Errorf("Incorrect result. Expect %+v, got %+v", testCase.expectedRepo, result)
		}
	}
}

func Test_DeleteTask(t *testing.T) {
	cfg, err := config.InitConfig("../../config.yaml")
	require.NoError(t, err)

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_repository.NewMockRepository(controller)

	ctx := context.Background()

	table := []struct {
		id primitive.ObjectID
	}{
		{
			id: primitive.NewObjectID(),
		},
	}

	for _, testCase := range table {
		mockRepo.EXPECT().DeleteTask(ctx, testCase.id).Return(nil).Times(1)

		service := New(mockRepo, cfg)

		err = service.DeleteTask(ctx, testCase.id)
		require.NoError(t, err)
	}
}
