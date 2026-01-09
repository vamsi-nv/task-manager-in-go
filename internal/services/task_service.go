package services

import (
	"context"
	"strconv"
	"time"

	"task-manager/internal/models"
	"task-manager/internal/repository"
	"task-manager/internal/utils"
	"task-manager/internal/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskService struct {
	Repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{
		Repo: repo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.CreateTaskRequest) (*models.Task, error) {
	userObjId, err := primitive.ObjectIDFromHex(ctx.Value("user_id").(string))
	if err != nil {
		return nil, utils.Unauthorized("Unauthorized access", nil)
	}

	var due time.Time
	if task.DueDate != "" {
		due, _ = time.Parse(time.RFC3339, task.DueDate)
	}
	newTask := &models.Task{
		UserID:      userObjId,
		Title:       task.Title,
		Description: task.Description,
		Category:    task.Category,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     due,
	}
	err = s.Repo.CreateTask(ctx, newTask)
	if err != nil {
		return nil, utils.Internal("Error creating task", nil)
	}

	return newTask, nil
}

func (s *TaskService) GetTasks(ctx context.Context, filters map[string]string) ([]models.Task, error) {

	userObjId, err := primitive.ObjectIDFromHex(ctx.Value("user_id").(string))
	if err != nil {
		return nil, utils.Unauthorized("Unauthorized access", nil)
	}
	filter := bson.M{"user_id": userObjId}

	if v, ok := filters["category"]; ok && v != "" {
		filter["category"] = v
	}
	if v, ok := filters["status"]; ok && v != "" {
		filter["status"] = v
	}

	if v, ok := filters["search"]; ok && v != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": v, "$options": "i"}},
			{"description": bson.M{"$regex": v, "$options": "i"}},
		}
	}

	sort := bson.D{}
	if v, ok := filters["sort"]; ok && v != "" {
		order := 1
		if v2, ok2 := filters["order"]; ok2 && v2 == "desc" {
			order = -1
		}
		sort = append(sort, bson.E{Key: v, Value: order})
	}

	limit := 20
	if v, ok := filters["limit"]; ok && v != "" {
		limit, _ = strconv.Atoi(v)
	}

	skip := 0
	if v, ok := filters["page"]; ok && v != "" {
		page, _ := strconv.Atoi(v)
		if page > 1 {
			skip = (page - 1) * limit
		}
	}

	tasks, err := s.Repo.GetTasks(ctx, filter, sort, limit, skip)
	if err != nil {
		return nil, utils.Internal("Error getting tasks", nil)
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id primitive.ObjectID, req models.UpdateTaskRequest) (*models.Task, error) {

	if !req.HasUpdates() {
		return nil, utils.BadRequest("no fields to update", nil)
	}

	err := validation.Validate.Struct(req)
	if err != nil {
		errs := utils.FormatValidationErrors(err)
		return nil, utils.BadRequest("Validation failed", errs)
	}

	task, err := s.Repo.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userObjId, err := primitive.ObjectIDFromHex(ctx.Value("user_id").(string))
	if err != nil {
		return nil, utils.Unauthorized("Unauthorized access", nil)
	}
	if task.UserID != userObjId {
		return nil, utils.Unauthorized("Unauthorized access", nil)
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}

	var due time.Time
	if req.DueDate != nil {
		due, _ = time.Parse(time.RFC3339, *req.DueDate)
		task.DueDate = due
	}

	updatedTask, err := s.Repo.UpdateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id primitive.ObjectID) error {

	task, err := s.Repo.GetTaskByID(ctx, id)
	if err != nil {
		return err
	}
	userObjId, err := primitive.ObjectIDFromHex(ctx.Value("user_id").(string))
	if err != nil {
		return utils.Unauthorized("Unauthorized access", nil)
	}
	if task.UserID != userObjId {
		return utils.Unauthorized("Unauthorized access", nil)
	}

	err = s.Repo.DeleteTask(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
