package projects

import (
	"context"
	"task-management/internal/entity"
	basic_service "task-management/internal/service/_basic_service"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) GetProjectsWithStats(ctx context.Context, filter Filter) ([]List, error) {
	return s.repo.GetProjectsWithStats(ctx, filter)
}

func (s Service) GetProjectsCount(ctx context.Context, filter Filter) (int, error) {
	return s.repo.GetProjectsCount(ctx, filter)
}

func (s Service) GetById(ctx context.Context, id int) (Detail, error) {
	return s.repo.GetById(ctx, id)
}

func (s Service) Create(ctx context.Context, data Create) (entity.Projects, error) {
	return s.repo.Create(ctx, data)
}

func (s Service) Update(ctx context.Context, data Update) (entity.Projects, error) {
	return s.repo.Update(ctx, data)
}

func (s Service) Delete(ctx context.Context, data basic_service.Delete) error {
	return s.repo.Delete(ctx, data)
}
