package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type TaskUpdate struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
	Status  *string `json:"status,omitempty"`
}

func (t *TaskUpdate) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(t.Title, validation.When(t.Title != nil, validation.Required.Error("Title cannot be empty"))),
		validation.Field(t.Status, validation.When(t.Status != nil, validation.In(StatusCompleted, StatusIncomplete, StatusInProgress).Error("Invalid status"))),
	)
}
