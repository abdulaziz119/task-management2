package tasks

import (
	"context"
	"task-management2/internal/entity"
	basic_repo "task-management2/internal/repository/postgres/_basic_repo"
	"task-management2/internal/repository/postgres/tasks"
)

type Repository interface {
	GetAll(ctx context.Context, filter tasks.Filter) ([]entity.Tasks, int, error)
	GetTaskStats(ctx context.Context, filter tasks.Filter) (tasks.TaskStats, error)
	GetById(ctx context.Context, id int) (entity.Tasks, error)
	Create(ctx context.Context, data tasks.Create) (entity.Tasks, error)
	Update(ctx context.Context, data tasks.Update) (entity.Tasks, error)
	Delete(ctx context.Context, data basic_repo.Delete) error
}
