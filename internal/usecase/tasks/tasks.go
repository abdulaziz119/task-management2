package tasks

import (
	"context"
	"task-management/internal/entity"
	basic_service "task-management/internal/service/_basic_service"
	"task-management/internal/service/tasks"
)

type TasksUseCase struct {
	service *tasks.Service
}

func NewUseCase(service *tasks.Service) *TasksUseCase {
	return &TasksUseCase{service: service}
}

func (uc TasksUseCase) GetAll(ctx context.Context, filter tasks.Filter) ([]entity.Tasks, tasks.TaskStats, int, error) {
	taskList, count, err := uc.service.GetAll(ctx, filter)
	if err != nil {
		return nil, tasks.TaskStats{}, 0, err
	}

	stats, err := uc.service.GetTaskStats(ctx, filter)
	if err != nil {
		return nil, tasks.TaskStats{}, 0, err
	}

	return taskList, stats, count, nil
}

func (uc TasksUseCase) TaskGetList(ctx context.Context, filter tasks.Filter) ([]tasks.List, tasks.TaskStats, int, error) {
	data, taskStats, count, err := uc.GetAll(ctx, filter)
	if err != nil {
		return nil, tasks.TaskStats{}, 0, err
	}

	var list []tasks.List
	for _, d := range data {
		var item tasks.List
		item.Id = d.Id
		if d.Name != nil {
			item.Name = *d.Name
		}
		if d.Description != nil {
			item.Description = *d.Description
		}
		if d.ProjectId != nil {
			item.ProjectId = *d.ProjectId
		}
		if d.Status != nil {
			item.Status = *d.Status
		}
		if d.Priority != nil {
			item.Priority = *d.Priority
		}
		if d.DueDate != nil {
			item.DueDate = *d.DueDate
		}
		if d.AssignedTo != nil {
			item.AssignedTo = *d.AssignedTo
		}

		list = append(list, item)
	}

	return list, taskStats, count, nil
}

func (uc TasksUseCase) TaskGetDetail(ctx context.Context, id int) (tasks.Detail, error) {
	data, err := uc.service.GetById(ctx, id)
	if err != nil {
		return tasks.Detail{}, err
	}

	var detail tasks.Detail

	detail.Id = data.Id
	detail.Name = data.Name
	detail.Description = data.Description
	detail.AssignedTo = data.AssignedTo
	detail.Status = data.Status
	detail.Priority = data.Priority
	detail.DueDate = data.DueDate

	return detail, nil
}

func (uc TasksUseCase) TaskCreate(ctx context.Context, data tasks.Create) (tasks.Detail, error) {
	task, err := uc.service.Create(ctx, data)
	if err != nil {
		return tasks.Detail{}, err
	}

	var detail tasks.Detail

	detail.Id = task.Id
	detail.Name = task.Name
	detail.Description = task.Description
	detail.AssignedTo = task.AssignedTo
	detail.Status = task.Status
	detail.Priority = task.Priority
	detail.DueDate = task.DueDate

	return detail, nil
}

func (uc TasksUseCase) TaskUpdate(ctx context.Context, data tasks.Update) (tasks.Detail, error) {
	task, err := uc.service.Update(ctx, data)
	if err != nil {
		return tasks.Detail{}, err
	}

	var detail tasks.Detail

	detail.Id = task.Id
	detail.Name = task.Name
	detail.Description = task.Description
	detail.AssignedTo = task.AssignedTo
	detail.Status = task.Status
	detail.Priority = task.Priority
	detail.DueDate = task.DueDate

	return detail, nil
}

func (uc TasksUseCase) TaskDelete(ctx context.Context, data basic_service.Delete) error {
	return uc.service.Delete(ctx, data)
}
