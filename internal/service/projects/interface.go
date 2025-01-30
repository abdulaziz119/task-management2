package projects

import (
	"context"
	"task-management/internal/entity"
	basic_service "task-management/internal/service/_basic_service"
)

type Repository interface {
	GetProjectsWithStats(ctx context.Context, filter Filter) ([]List, error)
	GetProjectsCount(ctx context.Context, filter Filter) (int, error)
	GetById(ctx context.Context, id int) (Detail, error)
	Create(ctx context.Context, data Create) (entity.Projects, error)
	Update(ctx context.Context, data Update) (entity.Projects, error)
	Delete(ctx context.Context, data basic_service.Delete) error
}
