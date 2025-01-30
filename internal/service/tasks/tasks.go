package tasks

import (
	"context"
	"errors"
	"task-management/internal/entity"
	basic_service "task-management/internal/service/_basic_service"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s Service) GetAll(ctx context.Context, filter Filter) ([]entity.Tasks, int, error) {
	return s.repo.GetAll(ctx, filter)
}

func (s Service) GetById(ctx context.Context, id int) (entity.Tasks, error) {
	return s.repo.GetById(ctx, id)
}

func (s Service) Create(ctx context.Context, data Create) (entity.Tasks, error) {
	if data.ProjectId == nil {
		return entity.Tasks{}, errors.New("project_id is required!")
	}
	if data.Name == nil {
		return entity.Tasks{}, errors.New("name is required!")
	}
	if data.Description == nil {
		return entity.Tasks{}, errors.New("description is required!")
	}
	if data.AssignedTo == nil {
		return entity.Tasks{}, errors.New("assigne_to is required!")
	}
	if data.Status == nil {
		return entity.Tasks{}, errors.New("status is required!")
	}
	if data.Priority == nil {
		return entity.Tasks{}, errors.New("Priority is required!")
	}

	return s.repo.Create(ctx, data)
}

func (s Service) Update(ctx context.Context, data Update) (entity.Tasks, error) {
	return s.repo.Update(ctx, data)
}

func (s Service) Delete(ctx context.Context, data basic_service.Delete) error {
	return s.repo.Delete(ctx, data)
}

func (s Service) GetTaskStats(ctx context.Context, filter Filter) (TaskStats, error) {
	return s.repo.GetTaskStats(ctx, filter)
}
