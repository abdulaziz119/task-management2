package users

import (
	"context"
	"task-management2/internal/entity"
	basic_repo "task-management2/internal/repository/postgres/_basic_repo"
	"task-management2/internal/repository/postgres/users"
)

type Repository interface {
	GetAll(ctx context.Context, filter users.Filter) ([]users.List, int, error)
	GetById(ctx context.Context, id int) (users.Detail, error)
	Create(ctx context.Context, data users.Create) (entity.User, error)
	Update(ctx context.Context, data users.Update) (entity.User, error)
	Delete(ctx context.Context, data basic_repo.Delete) error
}
