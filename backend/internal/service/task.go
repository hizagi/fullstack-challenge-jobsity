package service

import (
	"context"
	"time"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/api/generated"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/storage/model"
	"go.mongodb.org/mongo-driver/bson"
)

const defaultLimitValue int64 = 10

type taskRepository interface {
	CreateTask(ctx context.Context, task model.Task) (string, error)
	UpdateTask(ctx context.Context, id string, taskUpdate bson.M) error
	DeleteTask(ctx context.Context, id string) error
	GetTask(ctx context.Context, id string) (*model.Task, error)
	ListTasks(ctx context.Context, cursor string, limit int64) ([]model.Task, string, error)
}

type TaskService struct {
	taskRepository taskRepository
}

func NewTaskService(taskRepository taskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}

func (s *TaskService) CreateTask(ctx context.Context, createTask generated.CreateTask) (string, error) {
	newTask := model.Task{
		Title:     createTask.Title,
		Content:   createTask.Content,
		Status:    string(generated.TaskStatusIncomplete),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.taskRepository.CreateTask(ctx, newTask)
}

func (s *TaskService) UpdateTask(ctx context.Context, id string, updateTask generated.UpdateTask) error {
	update := make(map[string]interface{})

	if updateTask.Title != nil {
		update["title"] = *updateTask.Title
	}
	if updateTask.Content != nil {
		update["content"] = *updateTask.Content
	}
	if updateTask.Status != nil {
		update["status"] = *updateTask.Status
	}
	update["updateAt"] = time.Now()

	return s.taskRepository.UpdateTask(ctx, id, update)
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	return s.taskRepository.DeleteTask(ctx, id)
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*generated.Task, error) {
	model, err := s.taskRepository.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	hexID := model.ID.Hex()

	return &generated.Task{
		ID:        &hexID,
		Title:     &model.Title,
		Content:   &model.Content,
		Status:    s.convertStatus(model.Status),
		CreatedAt: &model.CreatedAt,
		UpdatedAt: &model.UpdatedAt,
	}, nil
}

func (s *TaskService) convertStatus(status string) *generated.TaskStatus {
	switch status {
	case string(generated.TaskStatusComplete):
		val := generated.TaskStatusComplete
		return &val
	case string(generated.TaskStatusInProgress):
		val := generated.TaskStatusInProgress
		return &val
	case string(generated.TaskStatusIncomplete):
		val := generated.TaskStatusIncomplete
		return &val
	default:
		return nil
	}
}

func (s *TaskService) ListTasks(ctx context.Context, cursor *string, limit *int64) ([]generated.Task, string, error) {
	cursorValue := ""
	limitValue := defaultLimitValue

	if cursor != nil {
		cursorValue = *cursor
	}

	if limit != nil {
		limitValue = *limit
	}

	models, nextCursor, err := s.taskRepository.ListTasks(ctx, cursorValue, limitValue)
	if err != nil || len(models) == 0 {
		return nil, "", err
	}

	var tasks []generated.Task
	for _, model := range models {
		hexID := model.ID.Hex()
		tasks = append(tasks, generated.Task{
			ID:        &hexID,
			Title:     &model.Title,
			Content:   &model.Content,
			Status:    s.convertStatus(model.Status),
			CreatedAt: &model.CreatedAt,
			UpdatedAt: &model.UpdatedAt,
		})
	}

	return tasks, nextCursor, nil
}
