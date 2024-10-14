package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/api/generated"
)

func ValidateCreateTask(createTask *generated.CreateTask) error {
	return validation.ValidateStruct(
		createTask,
		validation.Field(&createTask.Title, validation.Required.Error("Title is required")),
	)
}

func ValidateUpdateTask(updateTask *generated.UpdateTask) error {
	return validation.ValidateStruct(
		updateTask,
		validation.Field(&updateTask.Title, validation.Required.Error("Title is required")),
		validation.Field(&updateTask.Status, validation.Required.Error("Status is required"), validation.In(generated.TaskStatusComplete, generated.TaskStatusIncomplete, generated.TaskStatusInProgress).Error("Invalid status")),
	)
}
