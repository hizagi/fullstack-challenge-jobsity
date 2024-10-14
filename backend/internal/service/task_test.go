package service

import (
	"context"
	"testing"
	"time"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/api/generated"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/domain"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/storage/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestTaskService(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	mockTimeProvider := domain.NewMockTimeProvider(ctrl)
	mockTaskRepository := NewMockTaskRepository(ctrl)

	taskService := NewTaskService(mockTaskRepository, mockTimeProvider)

	defaultDate := time.Date(2024, 10, 14, 01, 05, 00, 0, time.UTC)
	testTitle := "Test Title"
	testContent := "Test Content"
	testIncomplete := generated.TaskStatusIncomplete
	testInProgressStatus := generated.UpdateTaskStatusInProgress

	t.Run("create task successfully", func(t *testing.T) {
		task := generated.CreateTask{
			Title:   testTitle,
			Content: testContent,
		}

		taskModel := model.Task{
			Title:     testTitle,
			Content:   testContent,
			Status:    string(generated.TaskStatusIncomplete),
			CreatedAt: defaultDate,
			UpdatedAt: defaultDate,
		}

		mockTimeProvider.EXPECT().Now().Return(defaultDate).Times(2)
		mockTaskRepository.EXPECT().CreateTask(gomock.Any(), taskModel).Return("1", nil).Times(1)

		taskID, err := taskService.CreateTask(ctx, task)
		assert.Nil(t, err)
		assert.Equal(t, "1", taskID)
	})

	t.Run("create task failure", func(t *testing.T) {
		task := generated.CreateTask{
			Title:   testTitle,
			Content: testContent,
		}

		taskModel := model.Task{
			Title:     testTitle,
			Content:   testContent,
			Status:    string(generated.TaskStatusIncomplete),
			CreatedAt: defaultDate,
			UpdatedAt: defaultDate,
		}

		mockTimeProvider.EXPECT().Now().Return(defaultDate).Times(2)
		mockTaskRepository.EXPECT().CreateTask(gomock.Any(), taskModel).Return("", assert.AnError).Times(1)

		taskID, err := taskService.CreateTask(ctx, task)
		assert.Equal(t, assert.AnError, err)
		assert.Equal(t, "", taskID)
	})

	t.Run("update task title successfully", func(t *testing.T) {
		task := generated.UpdateTask{
			Title: &testTitle,
		}

		taskID := "6526d630bac500ba72368e32"

		updateMap := map[string]interface{}{
			"title":     testTitle,
			"updatedAt": defaultDate,
		}

		mockTimeProvider.EXPECT().Now().Return(defaultDate).Times(1)
		mockTaskRepository.EXPECT().UpdateTask(gomock.Any(), taskID, updateMap).Return(nil).Times(1)

		err := taskService.UpdateTask(ctx, taskID, task)
		assert.Nil(t, err)
	})

	t.Run("update task title and status successfully", func(t *testing.T) {
		task := generated.UpdateTask{
			Title:  &testTitle,
			Status: &testInProgressStatus,
		}

		taskID := "6526d630bac500ba72368e32"

		updateMap := map[string]interface{}{
			"title":     testTitle,
			"status":    testInProgressStatus,
			"updatedAt": defaultDate,
		}

		mockTimeProvider.EXPECT().Now().Return(defaultDate).Times(1)
		mockTaskRepository.EXPECT().UpdateTask(gomock.Any(), taskID, updateMap).Return(nil).Times(1)

		err := taskService.UpdateTask(ctx, taskID, task)
		assert.Nil(t, err)
	})

	t.Run("failure updating task title and status", func(t *testing.T) {
		task := generated.UpdateTask{
			Title:  &testTitle,
			Status: &testInProgressStatus,
		}

		taskID := "6526d630bac500ba72368e32"

		updateMap := map[string]interface{}{
			"title":     testTitle,
			"status":    testInProgressStatus,
			"updatedAt": defaultDate,
		}

		mockTimeProvider.EXPECT().Now().Return(defaultDate).Times(1)
		mockTaskRepository.EXPECT().UpdateTask(gomock.Any(), taskID, updateMap).Return(assert.AnError).Times(1)

		err := taskService.UpdateTask(ctx, taskID, task)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("delete task successfully", func(t *testing.T) {
		taskID := "6526d630bac500ba72368e32"

		mockTaskRepository.EXPECT().DeleteTask(gomock.Any(), taskID).Return(nil).Times(1)

		err := taskService.DeleteTask(ctx, taskID)
		assert.Nil(t, err)
	})

	t.Run("failure deleting task", func(t *testing.T) {
		taskID := "6526d630bac500ba72368e32"

		mockTaskRepository.EXPECT().DeleteTask(gomock.Any(), taskID).Return(assert.AnError).Times(1)

		err := taskService.DeleteTask(ctx, taskID)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("get task successfully", func(t *testing.T) {
		taskID := "6526d630bac500ba72368e32"

		task := generated.Task{
			ID:        &taskID,
			Title:     &testTitle,
			Content:   &testContent,
			Status:    &testIncomplete,
			CreatedAt: &defaultDate,
			UpdatedAt: &defaultDate,
		}

		objectID, err := primitive.ObjectIDFromHex(taskID)
		require.NoErrorf(t, err, "failure decoding the objectId")

		taskModel := model.Task{
			ID:        objectID,
			Title:     testTitle,
			Content:   testContent,
			Status:    string(generated.TaskStatusIncomplete),
			CreatedAt: defaultDate,
			UpdatedAt: defaultDate,
		}

		mockTaskRepository.EXPECT().GetTask(gomock.Any(), taskID).Return(&taskModel, nil).Times(1)

		taskFound, err := taskService.GetTask(ctx, taskID)
		assert.Nil(t, err)
		assert.Equal(t, &task, taskFound)
	})

	t.Run("failed to get task", func(t *testing.T) {
		taskID := "6526d630bac500ba72368e32"

		mockTaskRepository.EXPECT().GetTask(gomock.Any(), taskID).Return(nil, assert.AnError).Times(1)

		taskFound, err := taskService.GetTask(ctx, taskID)
		assert.Equal(t, assert.AnError, err)
		assert.Nil(t, taskFound)
	})

	t.Run("list task successfully", func(t *testing.T) {
		taskID1 := "6526d630bac500ba72368e32"
		taskID2 := "6526d68fbac500ba72368e33"
		taskID3 := "6526d68fbac500ba72368e34"
		newCursorObjectID := "6526d68fbac500ba72368e35"

		objectID1, err := primitive.ObjectIDFromHex(taskID1)
		require.NoErrorf(t, err, "failure decoding the objectId")
		objectID2, err := primitive.ObjectIDFromHex(taskID2)
		require.NoErrorf(t, err, "failure decoding the objectId")
		objectID3, err := primitive.ObjectIDFromHex(taskID3)
		require.NoErrorf(t, err, "failure decoding the objectId")

		cursor := ""
		var limit int64 = 20

		tasks := []generated.Task{
			generated.Task{
				ID:        &taskID1,
				Title:     &testTitle,
				Content:   &testContent,
				Status:    &testIncomplete,
				CreatedAt: &defaultDate,
				UpdatedAt: &defaultDate,
			},
			generated.Task{
				ID:        &taskID2,
				Title:     &testTitle,
				Content:   &testContent,
				Status:    &testIncomplete,
				CreatedAt: &defaultDate,
				UpdatedAt: &defaultDate,
			},
			generated.Task{
				ID:        &taskID3,
				Title:     &testTitle,
				Content:   &testContent,
				Status:    &testIncomplete,
				CreatedAt: &defaultDate,
				UpdatedAt: &defaultDate,
			},
		}

		modelTasks := []model.Task{
			model.Task{
				ID:        objectID1,
				Title:     testTitle,
				Content:   testContent,
				Status:    string(generated.TaskStatusIncomplete),
				CreatedAt: defaultDate,
				UpdatedAt: defaultDate,
			},
			model.Task{
				ID:        objectID2,
				Title:     testTitle,
				Content:   testContent,
				Status:    string(generated.TaskStatusIncomplete),
				CreatedAt: defaultDate,
				UpdatedAt: defaultDate,
			},
			model.Task{
				ID:        objectID3,
				Title:     testTitle,
				Content:   testContent,
				Status:    string(generated.TaskStatusIncomplete),
				CreatedAt: defaultDate,
				UpdatedAt: defaultDate,
			},
		}

		mockTaskRepository.EXPECT().ListTasks(gomock.Any(), cursor, limit).Return(modelTasks, newCursorObjectID, nil).Times(1)

		tasksListed, newCursor, err := taskService.ListTasks(ctx, &cursor, &limit)
		assert.Nil(t, err)
		assert.Equal(t, newCursorObjectID, newCursor)
		assert.Equal(t, tasks, tasksListed)
	})

	t.Run("failure listing tasks", func(t *testing.T) {
		cursor := ""
		var limit int64 = 20

		mockTaskRepository.EXPECT().ListTasks(gomock.Any(), cursor, limit).Return(nil, "", assert.AnError).Times(1)

		tasksListed, newCursor, err := taskService.ListTasks(ctx, &cursor, &limit)
		assert.Equal(t, assert.AnError, err)
		assert.Equal(t, "", newCursor)
		assert.Nil(t, tasksListed)
	})
}
