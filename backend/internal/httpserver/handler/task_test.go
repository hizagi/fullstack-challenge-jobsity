package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/api/generated"
	"github.com/stretchr/testify/assert"
)

func TestTaskHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTaskService(ctrl)
	handler := NewTaskHandler(mockService)

	testTitle := "Test Title"
	testContent := "Test Content"
	testIncomplete := generated.TaskStatusIncomplete
	testInProgressStatus := generated.UpdateTaskStatusInProgress

	taskID1 := "6526d630bac500ba72368e32"
	taskID2 := "6526d68fbac500ba72368e33"
	taskID3 := "6526d68fbac500ba72368e34"

	defaultDate := time.Date(2024, 10, 14, 01, 05, 00, 0, time.UTC)

	t.Run("Valid Task Creation", func(t *testing.T) {
		// Prepare valid request
		validTask := generated.CreateTask{
			Title:   "Test Task",
			Content: "This is a test task",
		}
		mockService.EXPECT().CreateTask(gomock.Any(), validTask).Return("123", nil)

		body, _ := json.Marshal(validTask)
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		// Invoke the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "/tasks/123", rec.Header().Get("Location"))
	})

	t.Run("Invalid Task Creation", func(t *testing.T) {
		// Prepare invalid request
		invalidTask := generated.CreateTask{}

		body, _ := json.Marshal(invalidTask)
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		// Invoke the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Get Task Success", func(t *testing.T) {
		// Prepare expected task
		expectedTask := &generated.Task{
			ID:      &taskID1,
			Title:   &testTitle,
			Content: &testContent,
			Status:  &testIncomplete,
		}
		mockService.EXPECT().GetTask(gomock.Any(), "123").Return(expectedTask, nil)

		req := httptest.NewRequest(http.MethodGet, "/tasks/123", nil)
		rec := httptest.NewRecorder()

		// Invoke the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusOK, rec.Code)

		var task generated.Task
		err := json.NewDecoder(rec.Body).Decode(&task)
		assert.NoError(t, err)
		assert.Equal(t, taskID1, *task.ID)
		assert.Equal(t, testTitle, *task.Title)
	})

	t.Run("Get Task Failure", func(t *testing.T) {
		// Prepare error case
		mockService.EXPECT().GetTask(gomock.Any(), "123").Return(nil, errors.New("task not found"))

		req := httptest.NewRequest(http.MethodGet, "/tasks/123", nil)
		rec := httptest.NewRecorder()

		// Invoke the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
	t.Run("Delete Task Success", func(t *testing.T) {
		mockService.EXPECT().DeleteTask(gomock.Any(), "123").Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/tasks/123", nil)
		rec := httptest.NewRecorder()

		// Invoke the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Delete Task Failure", func(t *testing.T) {
		mockService.EXPECT().DeleteTask(gomock.Any(), "123").Return(errors.New("task not found"))

		req := httptest.NewRequest(http.MethodDelete, "/tasks/123", nil)
		rec := httptest.NewRecorder()

		// Invoke the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Update Task Success", func(t *testing.T) {
		// Prepare update request
		updateReq := generated.UpdateTask{
			Title:   &testTitle,
			Content: &testContent,
			Status:  &testInProgressStatus,
		}
		mockService.EXPECT().UpdateTask(gomock.Any(), "123", updateReq).Return(nil)

		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPatch, "/tasks/123", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		// Invoke the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Update Task Failure", func(t *testing.T) {
		// Prepare update request with error
		updateReq := generated.UpdateTask{
			Title:   &testTitle,
			Content: &testContent,
			Status:  &testInProgressStatus,
		}
		mockService.EXPECT().UpdateTask(gomock.Any(), "123", updateReq).Return(errors.New("task not found"))

		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPatch, "/tasks/123", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		// Invoke the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Update Invalid Task", func(t *testing.T) {
		invalidTestTitle := ""

		// Prepare update request with error
		updateReq := generated.UpdateTask{
			Title:   &invalidTestTitle,
			Content: &testContent,
		}

		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPatch, "/tasks/123", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		// Invoke the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("List Tasks Success", func(t *testing.T) {
		expectedTasks := []generated.Task{
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
		nextCursor := "next-cursor"

		mockService.EXPECT().ListTasks(gomock.Any(), gomock.Nil(), gomock.Nil()).Return(expectedTasks, nextCursor, nil)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		rec := httptest.NewRecorder()

		// Invoke the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusOK, rec.Code)

		var response generated.ListTasksResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, expectedTasks, response.Tasks)
		assert.Equal(t, nextCursor, response.NextCursor)
	})

	t.Run("List Tasks Failure", func(t *testing.T) {
		mockService.EXPECT().ListTasks(gomock.Any(), gomock.Nil(), gomock.Nil()).Return(nil, "", errors.New("error listing tasks"))

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		rec := httptest.NewRecorder()

		// Invoke the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
