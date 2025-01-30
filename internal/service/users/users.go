package users

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
	return &Service{repo: repo}
}

func (s Service) GetAll(ctx context.Context, filter Filter) ([]List, int, error) {
	return s.repo.GetAll(ctx, filter)
}

func (s Service) GetById(ctx context.Context, id int) (Detail, error) {
	return s.repo.GetById(ctx, id)
}

func (s Service) Create(ctx context.Context, data Create) (entity.User, error) {
	if data.Email == nil {
		return entity.User{}, errors.New("email is required!")
	}
	if data.Role == nil {
		return entity.User{}, errors.New("role is required!")
	}
	if data.FullName == nil {
		return entity.User{}, errors.New("full name is required!")
	}
	if data.Password == nil {
		return entity.User{}, errors.New("password is required")
	}

	return s.repo.Create(ctx, data)
}

func (s Service) Update(ctx context.Context, data Update) (entity.User, error) {
	return s.repo.Update(ctx, data)
}

func (s Service) Delete(ctx context.Context, data basic_service.Delete) error {
	return s.repo.Delete(ctx, data)
}
