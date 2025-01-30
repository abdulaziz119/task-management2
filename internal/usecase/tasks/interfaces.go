package tasks

import (
	"context"
	basic_service "task-management/internal/service/_basic_service"
	"task-management/internal/service/tasks"
)

type UseCase interface {
	TaskGetList(ctx context.Context, filter tasks.Filter) ([]tasks.List, tasks.TaskStats, int, error)
	TaskGetDetail(ctx context.Context, id int) (tasks.Detail, error)
	TaskCreate(ctx context.Context, data tasks.Create) (tasks.Detail, error)
	TaskUpdate(ctx context.Context, data tasks.Update) (tasks.Detail, error)
	TaskDelete(ctx context.Context, data basic_service.Delete) error
}
