//go:generate mockgen -source=task.go -destination=task_mocks.go -package=handler -mock_names=taskService=MockTaskService

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/api/generated"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/domain"
)

type taskService interface {
	CreateTask(ctx context.Context, createTask generated.CreateTask) (string, error)
	UpdateTask(ctx context.Context, id string, updateTask generated.UpdateTask) error
	DeleteTask(ctx context.Context, id string) error
	GetTask(ctx context.Context, id string) (*generated.Task, error)
	ListTasks(ctx context.Context, cursor *string, limit *int64) ([]generated.Task, string, error)
}

type TaskHandler struct {
	taskService taskService
}

func NewTaskHandler(taskService taskService, middlewares ...func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()

	middlewares = append(middlewares, middleware.Recoverer)

	corsOptions := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Add the CORS middleware to the router
	r.Use(corsOptions)

	r.Use(middlewares...)

	taskHandler := &TaskHandler{
		taskService: taskService,
	}

	return generated.HandlerWithOptions(taskHandler, generated.ChiServerOptions{
		BaseRouter: r,
	})
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request, params generated.ListTasksParams) {
	tasks, nextCursor, err := h.taskService.ListTasks(r.Context(), params.Cursor, params.Limit)
	if err != nil {
		log.Printf("Error listing tasks: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(generated.ListTasksResponse{
		Tasks:      tasks,
		NextCursor: nextCursor,
	})
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req generated.CreateTask
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := domain.ValidateCreateTask(&req); err != nil {
		jsonErrorResponse(w, err)
		return
	}

	taskID, err := h.taskService.CreateTask(r.Context(), req)
	if err != nil {
		log.Printf("Error creating task: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/tasks/%s", taskID))
	w.WriteHeader(http.StatusCreated)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request, id string) {
	var req generated.UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := domain.ValidateUpdateTask(&req); err != nil {
		jsonErrorResponse(w, err)
		return
	}

	err := h.taskService.UpdateTask(r.Context(), id, req)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request, id string) {
	err := h.taskService.DeleteTask(r.Context(), id)
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request, id string) {
	task, err := h.taskService.GetTask(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching task: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func jsonErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	statusCode := http.StatusInternalServerError
	body := make(map[string]interface{})
	// Switch to handle different error types
	switch err := err.(type) {
	case validation.Errors:
		statusCode = http.StatusBadRequest
		body["error"] = err

	default:
		body["error"] = "Internal Error, please contact an administrator."
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
