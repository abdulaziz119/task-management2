package projects

import (
	"context"
	"task-management/internal/entity"
	basic_service "task-management/internal/service/_basic_service"
	"task-management/internal/service/projects"
)

type Repository interface {
	GetProjectsWithStats(ctx context.Context, filter projects.Filter) ([]projects.List, error)
	GetProjectsCount(ctx context.Context, filter projects.Filter) (int, error)
	GetById(ctx context.Context, id int) (projects.Detail, error)
	Create(ctx context.Context, data projects.Create) (entity.Projects, error)
	Update(ctx context.Context, data projects.Update) (entity.Projects, error)
	Delete(ctx context.Context, data basic_service.Delete) error
}

type UseCase struct {
	projects Repository
}

func NewUseCase(projects Repository) *UseCase {
	return &UseCase{projects}
}

func (uc UseCase) ProjectGetList(ctx context.Context, filter projects.Filter) ([]projects.List, int, error) {
	list, err := uc.projects.GetProjectsWithStats(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	count, err := uc.projects.GetProjectsCount(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (uc UseCase) ProjectGetDetail(ctx context.Context, id int) (projects.Detail, error) {
	detail, err := uc.projects.GetById(ctx, id)
	if err != nil {
		return projects.Detail{}, err
	}

	return detail, nil
}

func (uc UseCase) ProjectCreate(ctx context.Context, data projects.Create) (entity.Projects, error) {
	return uc.projects.Create(ctx, data)
}

func (uc UseCase) ProjectUpdate(ctx context.Context, data projects.Update) (entity.Projects, error) {
	return uc.projects.Update(ctx, data)
}

func (uc UseCase) ProjectDelete(ctx context.Context, data basic_service.Delete) error {
	return uc.projects.Delete(ctx, data)
}
