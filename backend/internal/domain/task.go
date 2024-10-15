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
		validation.Field(&updateTask.Title, validation.NilOrNotEmpty.Error("Title should not be empty")),
		validation.Field(&updateTask.Status,
			validation.When(updateTask.Status != nil, validation.In(generated.UpdateTaskStatusComplete, generated.UpdateTaskStatusInProgress, generated.UpdateTaskStatusIncomplete)).Else(validation.Nil)),
	)
}
