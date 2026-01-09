package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"task-manager/internal/models"
	"task-manager/internal/services"
	"task-manager/internal/utils"
	"task-manager/internal/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	Service *services.TaskService
}

func NewTaskHandler(s *services.TaskService) *TaskHandler {
	return &TaskHandler{Service: s}
}

// disallows fields that are not present in the struct
func DecodeStrict[T any](r io.Reader, v *T) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// create task
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) error {
	var task models.CreateTaskRequest

	err := DecodeStrict(r.Body, &task)
	if err != nil {
		return utils.BadRequest("Invalid JSON", nil)
	}

	err = validation.Validate.Struct(task)
	if err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.BadRequest("Validation Failed", errs)
	}

	created, err := h.Service.CreateTask(r.Context(), &task)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusCreated, "Task created", created)
	return nil
}

// get tasks, filter, sort, paginate
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) error {

	filters := map[string]string{
		"category": "",
		"status":   "",
		"limit":    "",
		"page":     "",
		"sort":     "",
		"order":    "",
		"search":   "",
	}

	for key := range filters {
		if val, ok := r.URL.Query()[key]; ok {
			filters[key] = val[0]
		}
	}

	tasks, err := h.Service.GetTasks(r.Context(), filters)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusOK, "Tasks", struct {
		Count int           `json:"count"`
		Tasks []models.Task `json:"tasks"`
	}{
		Count: len(tasks),
		Tasks: tasks,
	})
	return nil
}

// update task by id
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	if id == "" {
		return utils.BadRequest("Task Id is required to update the task", nil)
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.BadRequest("Invalid task id", nil)
	}

	var body models.UpdateTaskRequest
	err = DecodeStrict(r.Body, &body)
	if err != nil {
		return utils.BadRequest("Invalid JSON", nil)
	}

	task, err := h.Service.UpdateTask(r.Context(), objectId, body)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusOK, "Task updated", task)
	return nil
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	if id == "" {
		return utils.BadRequest("Task Id is required to delete the task", nil)
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.BadRequest("Invalid task id", nil)
	}

	err = h.Service.DeleteTask(r.Context(), objectId)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusOK, "Task Deleted", nil)
	return nil
}
