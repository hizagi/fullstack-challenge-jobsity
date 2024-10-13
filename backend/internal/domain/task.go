package domain

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	StatusCompleted  = "completed"
	StatusIncomplete = "incomplete"
	StatusInProgress = "in-progress"
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t *Task) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Title, validation.Required.Error("Title is required")),
		validation.Field(&t.Status, validation.Required.Error("Status is required"), validation.In(StatusCompleted, StatusIncomplete, StatusInProgress).Error("Invalid status")),
	)
}
