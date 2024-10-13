package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/api/generated"
)

type TaskHandler struct {
}

func NewTaskHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	taskHandler := &TaskHandler{}

	return generated.Handler(taskHandler)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request, params generated.ListTasksParams) {
	log.Printf("ListTasks called, params: %v\n", params)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateTask called")
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Add logic for updating a task
	log.Println("UpdateTask called")
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request, id string) {
	log.Println("DeleteTask called for ID:", id)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request, id string) {
	log.Println("GetTask called for ID:", id)
}
