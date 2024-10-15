//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/storage"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/storage/model"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/tests/fixtures"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/tests/integration"
	"go.mongodb.org/mongo-driver/bson"
)

func TestTaskRepository(t *testing.T) {
	ctx := context.Background()

	// Start the MongoDB test container and get the DBConfig and close function
	dbConfig, close, err := integration.StartMongoContainer()
	if err != nil {
		t.Fatalf("failed to start mongo container: %v", err)
	}
	defer close(ctx)

	// Initialize MongoDB storage using the DBConfig returned from the container
	mongoStorage, err := storage.NewMongoStorage(ctx, *dbConfig)
	if err != nil {
		t.Fatalf("failed to initialize MongoDB storage: %v", err)
	}

	fixtures.Execute(mongoStorage.GetDatabase(), fixtures.TaskFixtureMethod)

	taskRepo := NewTaskRepository(mongoStorage)

	t.Run("successfully get task", func(t *testing.T) {
		task, err := taskRepo.GetTask(context.Background(), fixtures.PredefinedObjectID)
		require.NoErrorf(t, err, "error getting task")

		assert.Equal(t, fixtures.PredefinedObjectID, task.ID.Hex())
	})

	t.Run("successfully create task", func(t *testing.T) {
		task := model.Task{
			Title:   "Test Task",
			Content: "This is a test task",
			Status:  "incomplete",
		}

		taskID, err := taskRepo.CreateTask(context.Background(), task)
		require.NoErrorf(t, err, "error creating task")

		taskFound, err := taskRepo.GetTask(ctx, taskID)
		require.NoErrorf(t, err, "failed fetching task created")

		assert.Equal(t, taskID, taskFound.ID.Hex())
	})

	t.Run("successfully update task", func(t *testing.T) {
		task := model.Task{
			Title:   "Test Task",
			Content: "This is a test task that will be updated",
			Status:  "incomplete",
		}

		taskID, err := taskRepo.CreateTask(context.Background(), task)
		require.NoErrorf(t, err, "error creating task")

		titleUpdatedText := "Test Task Updated"

		err = taskRepo.UpdateTask(ctx, taskID, bson.M{"title": titleUpdatedText})
		require.NoErrorf(t, err, "failed to update task")

		taskFound, err := taskRepo.GetTask(ctx, taskID)
		require.NoErrorf(t, err, "failed fetching task created")

		assert.Equal(t, taskID, taskFound.ID.Hex())
		assert.Equal(t, task.Status, taskFound.Status)
		assert.Equal(t, task.Content, taskFound.Content)
		assert.Equal(t, titleUpdatedText, taskFound.Title)
	})

	t.Run("successfully delete task", func(t *testing.T) {
		task := model.Task{
			Title:   "Test Task to be deleted",
			Content: "This is a test task that will be deleted",
			Status:  "incomplete",
		}

		taskID, err := taskRepo.CreateTask(context.Background(), task)
		require.NoErrorf(t, err, "error creating task")

		err = taskRepo.DeleteTask(ctx, taskID)
		require.NoErrorf(t, err, "failed to update task")

		taskFound, err := taskRepo.GetTask(ctx, taskID)
		assert.Nil(t, taskFound)
		assert.Nil(t, err)
	})

	t.Run("successfully list tasks", func(t *testing.T) {
		cursor := ""
		totalTasks := 0

		for i := 0; i < 6; i++ {
			tasks, newCursor, err := taskRepo.ListTasks(ctx, cursor, 10)
			require.NoErrorf(t, err, "failed to list tasks")

			totalTasks = totalTasks + len(tasks)
			cursor = newCursor
		}

		assert.GreaterOrEqual(t, totalTasks, 50)
	})
}
