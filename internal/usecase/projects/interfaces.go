package projects

import (
	"context"
	"task-management/internal/entity"
	basic_service "task-management/internal/service/_basic_service"
	"task-management/internal/service/projects"
)

type Projects interface {
	GetAll(ctx context.Context, filter projects.Filter) ([]projects.List, int, error)
	GetById(ctx context.Context, id int) (projects.Detail, error)
	Create(ctx context.Context, data projects.Create) (entity.Projects, error)
	Update(ctx context.Context, data projects.Update) (entity.Projects, error)
	Delete(ctx context.Context, data basic_service.Delete) error
}
