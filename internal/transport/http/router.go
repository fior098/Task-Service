package http

import (
    "test-task/internal/transport/http/handlers"

    "github.com/gorilla/mux"
)

func NewRouter(taskHandler *handlers.TaskHandler) *mux.Router {
    r := mux.NewRouter()

    api := r.PathPrefix("/api/v1").Subrouter()

    api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
    api.HandleFunc("/tasks", taskHandler.ListTasks).Methods("GET")
    api.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
    api.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
    api.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

    return r
}