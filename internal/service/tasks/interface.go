package tasks

import (
	"context"
	"task-management/internal/entity"
	basic_service "task-management/internal/service/_basic_service"
)

type Repository interface {
	GetAll(ctx context.Context, filter Filter) ([]entity.Tasks, int, error)
	GetTaskStats(ctx context.Context, filter Filter) (TaskStats, error)
	GetById(ctx context.Context, id int) (entity.Tasks, error)
	Create(ctx context.Context, data Create) (entity.Tasks, error)
	Update(ctx context.Context, data Update) (entity.Tasks, error)
	Delete(ctx context.Context, data basic_service.Delete) error
}
