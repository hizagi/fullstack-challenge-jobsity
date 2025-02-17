// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package generated

import (
	"time"
)

// Defines values for TaskStatus.
const (
	TaskStatusComplete   TaskStatus = "complete"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusIncomplete TaskStatus = "incomplete"
)

// Defines values for UpdateTaskStatus.
const (
	UpdateTaskStatusComplete   UpdateTaskStatus = "complete"
	UpdateTaskStatusInProgress UpdateTaskStatus = "in-progress"
	UpdateTaskStatusIncomplete UpdateTaskStatus = "incomplete"
)

// CreateTask defines model for CreateTask.
type CreateTask struct {
	Content string `json:"content,omitempty"`
	Title   string `json:"title"`
}

// ListTasksResponse defines model for ListTasksResponse.
type ListTasksResponse struct {
	// NextCursor Cursor for the next page of tasks
	NextCursor string `json:"nextCursor"`
	Tasks      []Task `json:"tasks"`
}

// Task defines model for Task.
type Task struct {
	Content *string `json:"content,omitempty"`

	// CreatedAt Task creation time
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// ID Unique identifier for the task
	ID     *string     `json:"id,omitempty"`
	Status *TaskStatus `json:"status,omitempty"`
	Title  *string     `json:"title,omitempty"`

	// UpdatedAt Task update time
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// TaskStatus defines model for Task.Status.
type TaskStatus string

// UpdateTask defines model for UpdateTask.
type UpdateTask struct {
	Content *string           `json:"content,omitempty"`
	Status  *UpdateTaskStatus `json:"status,omitempty"`
	Title   *string           `json:"title,omitempty"`
}

// UpdateTaskStatus defines model for UpdateTask.Status.
type UpdateTaskStatus string

// ListTasksParams defines parameters for ListTasks.
type ListTasksParams struct {
	// Limit Maximum number of tasks to return
	Limit *int64 `form:"limit,omitempty" json:"limit,omitempty"`

	// Cursor Cursor (ID of the last task from the previous page)
	Cursor *string `form:"cursor,omitempty" json:"cursor,omitempty"`
}

// CreateTaskJSONRequestBody defines body for CreateTask for application/json ContentType.
type CreateTaskJSONRequestBody = CreateTask

// UpdateTaskJSONRequestBody defines body for UpdateTask for application/json ContentType.
type UpdateTaskJSONRequestBody = UpdateTask
