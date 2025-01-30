package users

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	basic_service "task-management/internal/service/_basic_service"
	"task-management/internal/service/users"
)

type Repository interface {
	GetAllUsers(ctx context.Context, filter users.Filter) ([]users.User, int, error)
	GetTaskStats(ctx context.Context) (map[int64]users.TaskStats, error)
}

type UseCase struct {
	user User
	repo Repository
}

var validate = validator.New()

func NewUseCase(user User, repo Repository) *UseCase {
	return &UseCase{user, repo}
}

func (uc UseCase) GetAll(ctx context.Context, filter users.Filter) ([]users.List, int, error) {
	usersList, count, err := uc.repo.GetAllUsers(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	taskStats, err := uc.repo.GetTaskStats(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get task stats: %w", err)
	}

	var result []users.List
	for _, user := range usersList {
		if user.Id == nil {
			continue
		}
		stats := taskStats[*user.Id]
		result = append(result, users.List{
			Id:              user.Id,
			FullName:        user.FullName,
			Email:           user.Email,
			Role:            user.Role,
			PendingTasks:    stats.PendingTasks,
			InProgressTasks: stats.InProgressTasks,
			CompletedTasks:  stats.CompletedTasks,
			TaskCount:       stats.TaskCount,
		})
	}

	return result, count, nil
}

func (uc UseCase) AdminGetUserDetail(ctx context.Context, id int) (users.Detail, error) {
	return uc.user.GetById(ctx, id)
}

func (uc UseCase) AdminCreateUser(ctx context.Context, data users.Create) error {
	if err := validate.Struct(data); err != nil {
		return err
	}
	_, err := uc.user.Create(ctx, data)
	return err
}

func (uc UseCase) AdminUpdateUser(ctx context.Context, data users.Update) error {
	if err := validate.Struct(data); err != nil {
		return err
	}
	_, err := uc.user.Update(ctx, data)
	return err
}

func (uc UseCase) AdminDeleteUser(ctx context.Context, data basic_service.Delete) error {
	return uc.user.Delete(ctx, data)
}
