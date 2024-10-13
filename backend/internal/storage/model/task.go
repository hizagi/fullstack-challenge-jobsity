package model

import (
	"time"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func FromDomain(domainTask domain.Task) Task {
	var objectID primitive.ObjectID
	if domainTask.ID != "" {
		objectID, _ = primitive.ObjectIDFromHex(domainTask.ID)
	}

	return Task{
		ID:        objectID,
		Title:     domainTask.Title,
		Content:   domainTask.Content,
		Status:    domainTask.Status,
		CreatedAt: domainTask.CreatedAt,
		UpdatedAt: domainTask.UpdatedAt,
	}
}

func (t Task) ToDomain() domain.Task {
	return domain.Task{
		ID:        t.ID.Hex(),
		Title:     t.Title,
		Content:   t.Content,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

type TaskCollection []Task

func (c TaskCollection) ToDomain() []domain.Task {
	var domainTasks []domain.Task

	for _, task := range c {
		domainTasks = append(domainTasks, task.ToDomain())
	}

	return domainTasks
}

func SetFromDomain(taskUpdate domain.TaskUpdate) bson.M {
	update := make(bson.M)

	if taskUpdate.Title != nil {
		update["title"] = *taskUpdate.Title
	}
	if taskUpdate.Content != nil {
		update["content"] = *taskUpdate.Content
	}
	if taskUpdate.Status != nil {
		update["status"] = *taskUpdate.Status
	}

	return bson.M{"$set": update}
}
