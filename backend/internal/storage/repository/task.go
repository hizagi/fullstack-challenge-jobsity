package repository

import (
	"context"
	"fmt"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/domain"
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

func (r *TaskRepository) DeleteTask(ctx context.Context, id string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid objectid: %w", err)
	}

	result, err := r.mongoCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return "", err
	}

	if result.DeletedCount == 1 {
		return id, nil
	}

	return "", fmt.Errorf("unable to delete ID: %s", id)
}

func (r *TaskRepository) GetTask(ctx context.Context, id string) (*domain.Task, error) {
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

	domainTask := task.ToDomain()

	return &domainTask, nil
}

func (r *TaskRepository) ListTasks(ctx context.Context, cursor string, limit int64) ([]domain.Task, string, error) {
	var tasks model.TaskCollection

	filter := bson.M{}
	if cursor != "" {
		// Convert the cursor to ObjectID
		objectID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor: %w", err)
		}
		// Only fetch tasks with _id greater than the cursor
		filter["_id"] = bson.M{"$gt": objectID}
	}

	// Define find options with a limit
	findOptions := options.Find().SetLimit(limit).SetSort(bson.M{"_id": 1}) // Ascending order by _id

	// Find tasks
	cursorMongo, err := r.mongoCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, "", err
	}
	defer cursorMongo.Close(ctx)

	// Iterate through the cursor and decode each task into the result array
	var lastTaskID primitive.ObjectID
	for cursorMongo.Next(ctx) {
		var task model.Task
		err := cursorMongo.Decode(&task)
		if err != nil {
			return nil, "", err
		}
		tasks = append(tasks, task)
		lastTaskID = task.ID // Capture the last task's ID for the next cursor
	}

	// Check for cursor errors
	if err := cursorMongo.Err(); err != nil {
		return nil, "", err
	}

	// If there are no results, return empty list and no next cursor
	if len(tasks) == 0 {
		return nil, "", nil
	}

	// Convert lastTaskID to hex string for the next cursor
	nextCursor := lastTaskID.Hex()

	// Return the list of tasks and the next cursor
	return tasks.ToDomain(), nextCursor, nil
}
