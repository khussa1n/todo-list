package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/khussa1n/todo-list/internal/config"
	"github.com/khussa1n/todo-list/internal/custom_error"
	"github.com/khussa1n/todo-list/internal/entity"
	"github.com/khussa1n/todo-list/internal/entity/dto"
	mock_repository "github.com/khussa1n/todo-list/internal/repository/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func Test_CreateTask(t *testing.T) {
	id := primitive.NewObjectID()

	table := []struct {
		name            string
		dto             dto.TasksDTO
		taskRepo        entity.Tasks
		expectedRepo    entity.Tasks
		expectedSrvc    entity.Tasks
		expectedSrvcErr error
	}{
		{
			name:         "ok_ВЫХОДНОЙ",
			dto:          dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-05"},
			taskRepo:     entity.Tasks{Title: "ВЫХОДНОЙ - Купить", ActiveAt: "2023-08-05", Status: "active"},
			expectedRepo: entity.Tasks{ID: id, Title: "ВЫХОДНОЙ - Купить", ActiveAt: "2023-08-05", Status: "active"},
			expectedSrvc: entity.Tasks{ID: id, Title: "ВЫХОДНОЙ - Купить", ActiveAt: "2023-08-05", Status: "active"},
		},
		{
			name:         "ok",
			dto:          dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-04"},
			taskRepo:     entity.Tasks{Title: "Купить", ActiveAt: "2023-08-04", Status: "active"},
			expectedRepo: entity.Tasks{ID: id, Title: "Купить", ActiveAt: "2023-08-04", Status: "active"},
			expectedSrvc: entity.Tasks{ID: id, Title: "Купить", ActiveAt: "2023-08-04", Status: "active"},
		},
		{
			name:            "activeAt invalid format",
			dto:             dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-32"},
			expectedSrvcErr: custom_error.ErrInvalidActiveAtFormat,
		},
		{
			name: "more than 200 char",
			dto: dto.TasksDTO{Title: "Купитьasdfsasddddddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd", ActiveAt: "2023-08-05"},
			expectedSrvcErr: custom_error.ErrMessageTooLong,
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			cfg, err := config.InitConfig("../../config.yaml")
			require.NoError(t, err)

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockRepo := mock_repository.NewMockRepository(controller)

			ctx := context.Background()

			service := New(mockRepo, cfg)

			switch testCase.name {
			case "ok", "ok ВЫХОДНОЙ":
				mockRepo.EXPECT().CreateTask(ctx, &testCase.taskRepo).Return(&testCase.expectedRepo, nil).Times(1)

				result, err := service.CreateTask(ctx, &testCase.dto)
				require.NoError(t, err)
				require.NotEmpty(t, result)

				require.Equal(t, *result, testCase.expectedSrvc)
				break
			case "activeAt invalid format", "more than 200 char":
				mockRepo.EXPECT().CreateTask(ctx, &testCase.taskRepo).Return(&testCase.expectedRepo, nil).Times(0)

				_, err = service.CreateTask(ctx, &testCase.dto)
				require.NotEmpty(t, err)

				require.Equal(t, err, testCase.expectedSrvcErr)
				break
			}
		})
	}
}

func Test_UpdateTask(t *testing.T) {
	id := primitive.NewObjectID()

	table := []struct {
		name            string
		dto             dto.TasksDTO
		taskRepo        primitive.ObjectID
		expectedRepo    entity.Tasks
		expectedRepo2   entity.Tasks
		expectedSrvcErr error
	}{
		{
			name:          "ok_ВЫХОДНОЙ",
			dto:           dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-05"},
			taskRepo:      id,
			expectedRepo:  entity.Tasks{ID: id, Title: "Купить", ActiveAt: "2023-08-05", Status: "active"},
			expectedRepo2: entity.Tasks{Title: "ВЫХОДНОЙ - Купить", ActiveAt: "2023-08-05", Status: "active"},
		},
		{
			name:          "ok",
			dto:           dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-04"},
			taskRepo:      id,
			expectedRepo:  entity.Tasks{ID: id, Title: "Купить", ActiveAt: "2023-08-04", Status: "active"},
			expectedRepo2: entity.Tasks{Title: "Купить", ActiveAt: "2023-08-04", Status: "active"},
		},
		{
			name: "more than 200 char",
			dto: dto.TasksDTO{Title: "Купитьasdfsasddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd" +
				"dddddddddddddddddddddddddddddddddddddddddddddddadsfdd", ActiveAt: "2023-08-04"},
			taskRepo:        id,
			expectedSrvcErr: custom_error.ErrMessageTooLong,
		},
		{
			name:            "activeAt invalid format",
			dto:             dto.TasksDTO{Title: "Купить", ActiveAt: "2023-08-70"},
			taskRepo:        id,
			expectedSrvcErr: custom_error.ErrInvalidActiveAtFormat,
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			cfg, err := config.InitConfig("../../config.yaml")
			require.NoError(t, err)

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockRepo := mock_repository.NewMockRepository(controller)

			ctx := context.Background()

			service := New(mockRepo, cfg)

			switch testCase.name {
			case "ok", "ok ВЫХОДНОЙ":
				mockRepo.EXPECT().GetTaskByID(ctx, testCase.taskRepo).Return(&testCase.expectedRepo, nil).Times(1)
				mockRepo.EXPECT().UpdateTask(ctx, &testCase.expectedRepo2, testCase.taskRepo).Return(nil).Times(1)

				err = service.UpdateTask(ctx, &testCase.dto, testCase.expectedRepo.ID)
				require.NoError(t, err)
				break
			case "activeAt invalid format", "more than 200 char":
				mockRepo.EXPECT().GetTaskByID(ctx, testCase.taskRepo).Return(&testCase.expectedRepo, nil).Times(0)
				mockRepo.EXPECT().UpdateTask(ctx, &testCase.expectedRepo2, testCase.taskRepo).Return(nil).Times(0)

				err = service.UpdateTask(ctx, &testCase.dto, testCase.expectedRepo.ID)
				require.Equal(t, testCase.expectedSrvcErr, err)
				break
			}
		})
	}
}

func Test_UpdateTaskStatus(t *testing.T) {
	table := []struct {
		name   string
		id     primitive.ObjectID
		status string
	}{
		{
			name:   "ok",
			id:     primitive.NewObjectID(),
			status: "done",
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			cfg, err := config.InitConfig("../../config.yaml")
			require.NoError(t, err)

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockRepo := mock_repository.NewMockRepository(controller)

			ctx := context.Background()

			mockRepo.EXPECT().UpdateTaskStatus(ctx, testCase.id, testCase.status).Return(nil).Times(1)

			service := New(mockRepo, cfg)

			err = service.UpdateTaskStatus(ctx, testCase.id, testCase.status)
			require.NoError(t, err)
		})
	}
}

func Test_GetAllTasks(t *testing.T) {
	table := []struct {
		name         string
		status       string
		expectedRepo []entity.Tasks
	}{
		{
			name:         "ok",
			status:       "active",
			expectedRepo: []entity.Tasks{{ID: primitive.NewObjectID(), Title: "Купить", ActiveAt: "2023-08-05", Status: "active"}},
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			cfg, err := config.InitConfig("../../config.yaml")
			require.NoError(t, err)

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockRepo := mock_repository.NewMockRepository(controller)

			ctx := context.Background()

			mockRepo.EXPECT().GetAllTasks(ctx, testCase.status).Return(testCase.expectedRepo, nil).Times(1)

			service := New(mockRepo, cfg)

			result, err := service.GetAllTasks(ctx, testCase.status)
			require.NoError(t, err)
			require.NotEmpty(t, result)
			require.Equal(t, testCase.expectedRepo[0], result[0])
		})
	}
}

func Test_DeleteTask(t *testing.T) {
	table := []struct {
		name string
		id   primitive.ObjectID
	}{
		{
			name: "ok",
			id:   primitive.NewObjectID(),
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			cfg, err := config.InitConfig("../../config.yaml")
			require.NoError(t, err)

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockRepo := mock_repository.NewMockRepository(controller)

			ctx := context.Background()

			mockRepo.EXPECT().DeleteTask(ctx, testCase.id).Return(nil).Times(1)

			service := New(mockRepo, cfg)

			err = service.DeleteTask(ctx, testCase.id)
			require.NoError(t, err)
		})
	}
}
