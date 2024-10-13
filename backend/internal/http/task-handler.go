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

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Add logic for creating a task (e.g., save to database)
	log.Println("CreateTask called")
	return
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Add logic for updating a task
	log.Println("UpdateTask called")
	return
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request, id string) {
	// Add logic for deleting a task (e.g., delete from database)
	log.Println("DeleteTask called for ID:", id)
	return
}
