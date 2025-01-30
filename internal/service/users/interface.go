package users

import (
	"context"
	"task-management/internal/entity"
	basic_service "task-management/internal/service/_basic_service"
)

type Repository interface {
	GetAll(ctx context.Context, filter Filter) ([]List, int, error)
	GetById(ctx context.Context, id int) (Detail, error)
	Create(ctx context.Context, data Create) (entity.User, error)
	Update(ctx context.Context, data Update) (entity.User, error)
	Delete(ctx context.Context, data basic_service.Delete) error
}
