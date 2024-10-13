package model

import (
	"time"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type Task struct {
	ID        string    `bson:"_id,omitempty"`
	Title     string    `bson:"title"`
	Content   string    `bson:"content"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func FromDomain(domainTask domain.Task) Task {
	return Task{
		ID:        domainTask.ID,
		Title:     domainTask.Title,
		Content:   domainTask.Content,
		Status:    domainTask.Status,
		CreatedAt: domainTask.CreatedAt,
		UpdatedAt: domainTask.UpdatedAt,
	}
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
