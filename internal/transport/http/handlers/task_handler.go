package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "test-task/internal/domain/task"
    taskUsecase "test-task/internal/usecase/task"

    "github.com/gorilla/mux"
)

type TaskHandler struct {
    service *taskUsecase.Service
}

func NewTaskHandler(service *taskUsecase.Service) *TaskHandler {
    return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var req CreateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    var recurrenceType *task.RecurrenceType
    if req.RecurrenceType != nil {
        rt := task.RecurrenceType(*req.RecurrenceType)
        recurrenceType = &rt
    }

    t := &task.Task{
        Title:                  req.Title,
        Description:            req.Description,
        Status:                 req.Status,
        ScheduledAt:            req.ScheduledAt,
        RecurrenceType:         recurrenceType,
        RecurrenceInterval:     req.RecurrenceInterval,
        RecurrenceDayOfMonth:   req.RecurrenceDayOfMonth,
        RecurrenceParity:       req.RecurrenceParity,
        RecurrenceSpecificDates: req.RecurrenceSpecificDates,
        RecurrenceEndDate:      req.RecurrenceEndDate,
    }

    created, err := h.service.CreateTask(r.Context(), t)
    if err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondJSON(w, http.StatusCreated, created)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondError(w, http.StatusBadRequest, "invalid task id")
        return
    }

    t, err := h.service.GetTask(r.Context(), id)
    if err != nil {
        if err == task.ErrTaskNotFound {
            respondError(w, http.StatusNotFound, "task not found")
            return
        }
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, t)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

    tasks, err := h.service.ListTasks(r.Context(), limit, offset)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondError(w, http.StatusBadRequest, "invalid task id")
        return
    }

    var req UpdateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    var recurrenceType *task.RecurrenceType
    if req.RecurrenceType != nil {
        rt := task.RecurrenceType(*req.RecurrenceType)
        recurrenceType = &rt
    }

    t := &task.Task{
        ID:                     id,
        Title:                  req.Title,
        Description:            req.Description,
        Status:                 req.Status,
        ScheduledAt:            req.ScheduledAt,
        RecurrenceType:         recurrenceType,
        RecurrenceInterval:     req.RecurrenceInterval,
        RecurrenceDayOfMonth:   req.RecurrenceDayOfMonth,
        RecurrenceParity:       req.RecurrenceParity,
        RecurrenceSpecificDates: req.RecurrenceSpecificDates,
        RecurrenceEndDate:      req.RecurrenceEndDate,
    }

    updated, err := h.service.UpdateTask(r.Context(), t)
    if err != nil {
        if err == task.ErrTaskNotFound {
            respondError(w, http.StatusNotFound, "task not found")
            return
        }
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, updated)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondError(w, http.StatusBadRequest, "invalid task id")
        return
    }

    if err := h.service.DeleteTask(r.Context(), id); err != nil {
        if err == task.ErrTaskNotFound {
            respondError(w, http.StatusNotFound, "task not found")
            return
        }
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
    respondJSON(w, status, ErrorResponse{Error: message})
}