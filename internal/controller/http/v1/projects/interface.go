package projects

import (
	"context"
	"task-management2/internal/entity"
	basic_repo "task-management2/internal/repository/postgres/_basic_repo"
	"task-management2/internal/repository/postgres/projects"
)

type Repository interface {
	GetProjectsWithStats(ctx context.Context, filter projects.Filter) ([]projects.List, error)
	GetProjectsCount(ctx context.Context, filter projects.Filter) (int, error)
	GetById(ctx context.Context, id int) (projects.Detail, error)
	Create(ctx context.Context, data projects.Create) (entity.Projects, error)
	Update(ctx context.Context, data projects.Update) (entity.Projects, error)
	Delete(ctx context.Context, data basic_repo.Delete) error
}
