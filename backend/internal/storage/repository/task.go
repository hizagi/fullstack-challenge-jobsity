package repository

import (
	"context"
	"fmt"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/storage"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/storage/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const taskCollection = "tasks"

type TaskRepository struct {
	mongoCollection *mongo.Collection
}

func NewTaskRepository(mongoStorage *storage.MongoStorage) *TaskRepository {
	return &TaskRepository{mongoCollection: mongoStorage.GetDatabase().Collection(taskCollection)}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task model.Task) (string, error) {
	result, err := r.mongoCollection.InsertOne(ctx, task)
	if err != nil {
		return "", err
	}

	if insertedID, ok := result.InsertedID.(primitive.ObjectID); ok {
		return insertedID.Hex(), nil
	}

	return "", fmt.Errorf("unable to convert inserted ID: %v", result.InsertedID)
}

func (r *TaskRepository) UpdateTask(ctx context.Context, id string, setMap primitive.M) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid objectid: %w", err)
	}

	result, err := r.mongoCollection.UpdateByID(ctx, objectID, setMap)
	if err != nil {
		return err
	}

	if result.UpsertedCount == 1 {
		return nil
	}

	return fmt.Errorf("unable to update ID: %s", id)
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid objectid: %w", err)
	}
	result, err := r.mongoCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 1 {
		return nil
	}

	return fmt.Errorf("unable to delete ID: %s", id)
}

func (r *TaskRepository) GetTask(ctx context.Context, id string) (*model.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid objectid: %w", err)
	}

	var task model.Task

	err = r.mongoCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	return &task, nil
}

func (r *TaskRepository) ListTasks(ctx context.Context, cursor string, limit int64) ([]model.Task, string, error) {
	var tasks []model.Task

	filter := bson.M{}
	if cursor != "" {
		objectID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor: %w", err)
		}

		filter["_id"] = bson.M{"$gt": objectID}
	}

	findOptions := options.Find().SetLimit(limit).SetSort(bson.M{"_id": 1})

	cursorMongo, err := r.mongoCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, "", err
	}
	defer cursorMongo.Close(ctx)

	var lastTaskID primitive.ObjectID
	for cursorMongo.Next(ctx) {
		var task model.Task
		err := cursorMongo.Decode(&task)
		if err != nil {
			return nil, "", err
		}
		tasks = append(tasks, task)
		lastTaskID = task.ID
	}

	if err := cursorMongo.Err(); err != nil {
		return nil, "", err
	}

	if len(tasks) == 0 {
		return nil, "", nil
	}

	nextCursor := lastTaskID.Hex()
	return tasks, nextCursor, nil
}
