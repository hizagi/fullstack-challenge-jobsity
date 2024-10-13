package repository

import (
	"context"
	"fmt"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/domain"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/storage"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/storage/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const taskCollection = "tasks"

type TaskRepository struct {
	mongoCollection *mongo.Collection
}

func NewTaskRepository(mongoStorage storage.MongoStorage) *TaskRepository {
	return &TaskRepository{mongoCollection: mongoStorage.GetDatabase().Collection(taskCollection)}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task domain.Task) (string, error) {
	result, err := r.mongoCollection.InsertOne(ctx, model.FromDomain(task))
	if err != nil {
		return "", err
	}

	if insertedID, ok := result.InsertedID.(primitive.ObjectID); ok {
		return insertedID.Hex(), nil
	}

	return "", fmt.Errorf("unable to convert inserted ID: %v", result.InsertedID)
}

func (r *TaskRepository) UpdateTask(ctx context.Context, id string, taskUpdate domain.TaskUpdate) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid objectid: %w", err)
	}

	result, err := r.mongoCollection.UpdateByID(ctx, objectID, model.SetFromDomain(taskUpdate))
	if err != nil {
		return "", err
	}

	if upsertedID, ok := result.UpsertedID.(primitive.ObjectID); ok {
		return upsertedID.Hex(), nil
	}

	return "", fmt.Errorf("unable to convert inserted ID: %v", result.UpsertedID)
}
