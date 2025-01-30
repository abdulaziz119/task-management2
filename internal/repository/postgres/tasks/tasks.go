package tasks

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"math"
	"task-management/internal/entity"
	basic_repo "task-management/internal/repository/postgres/_basic_repo"
	basic_service "task-management/internal/service/_basic_service"
	"task-management/internal/service/tasks"
	"time"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r Repository) GetAll(ctx context.Context, filter tasks.Filter) ([]entity.Tasks, int, error) {
	baseQuery := `
		WITH total_count AS (
			SELECT COUNT(*) as total
			FROM tasks t
			LEFT JOIN users u ON u.id = t.assigned_to
			WHERE t.deleted_at IS NULL
			%s
		)
		SELECT 
			t.id,
			t.project_id,
			t.name,
			t.description,
			t.assigned_to,
			t.status,
			t.priority,
			t.due_date,
			t.created_at,
			t.deleted_at,
			tc.total as total_count
		FROM tasks t
		LEFT JOIN users u ON u.id = t.assigned_to
		CROSS JOIN total_count tc
		WHERE t.deleted_at IS NULL
		%s
		ORDER BY t.id`

	projectFilter := ""
	if filter.ProjectId != nil {
		projectFilter = fmt.Sprintf("AND t.project_id = %d", *filter.ProjectId)
	}

	query := fmt.Sprintf(baseQuery, projectFilter, projectFilter)

	if filter.Offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *filter.Offset)
	}
	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying tasks: %v", err)
	}
	defer rows.Close()

	var result []entity.Tasks
	var totalCount int

	for rows.Next() {
		var task entity.Tasks

		err := rows.Scan(
			&task.Id,
			&task.ProjectId,
			&task.Name,
			&task.Description,
			&task.AssignedTo,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.DeletedAt,
			&totalCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning task row: %v", err)
		}

		result = append(result, task)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating task rows: %v", err)
	}

	return result, totalCount, nil
}

func (r Repository) GetTaskStats(ctx context.Context, filter tasks.Filter) (tasks.TaskStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_tasks,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_tasks,
			COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_tasks
		FROM tasks
		WHERE deleted_at IS NULL
	`

	if filter.ProjectId != nil {
		query += fmt.Sprintf(" AND project_id = %d", *filter.ProjectId)
	}

	var stats tasks.TaskStats
	var totalTasks, completedTasks, pendingTasks, inProgressTasks int

	err := r.QueryRowContext(ctx, query).Scan(
		&totalTasks,
		&completedTasks,
		&pendingTasks,
		&inProgressTasks,
	)
	if err != nil {
		return tasks.TaskStats{}, fmt.Errorf("error getting task stats: %v", err)
	}

	stats = tasks.TaskStats{
		TotalTasks:      totalTasks,
		CompletedTasks:  completedTasks,
		PendingTasks:    pendingTasks,
		InProgressTasks: inProgressTasks,
	}

	if stats.TotalTasks > 0 {
		completedProgress := float64(stats.CompletedTasks) * 100.0
		inProgressProgress := float64(stats.InProgressTasks) * 50.0
		totalPossibleProgress := float64(stats.TotalTasks) * 100.0

		stats.Progress = math.Round((completedProgress+inProgressProgress)/totalPossibleProgress*1000) / 10
	}

	return stats, nil
}

func (r Repository) GetById(ctx context.Context, id int) (entity.Tasks, error) {
	var detail entity.Tasks
	err := r.NewSelect().
		Model(&detail).
		Where("id = ? AND deleted_at IS NULL", id).
		Scan(ctx)

	if err != nil {
		return entity.Tasks{}, fmt.Errorf("error getting task: %v", err)
	}

	return detail, nil
}

func (r Repository) Create(ctx context.Context, data tasks.Create) (entity.Tasks, error) {
	var detail entity.Tasks

	const layout = "2006-01-02"
	_, err := time.Parse(layout, *data.DueDate)
	if err != nil {
		return entity.Tasks{}, fmt.Errorf("invalid DueDate format: %v", err)
	}

	detail.ProjectId = data.ProjectId
	detail.Name = data.Name
	detail.Description = data.Description
	detail.AssignedTo = data.AssignedTo
	detail.Status = data.Status
	detail.Priority = data.Priority
	detail.DueDate = data.DueDate

	_, err = r.NewInsert().Model(&detail).Exec(ctx)
	if err != nil {
		return entity.Tasks{}, fmt.Errorf("error creating task: %v", err)
	}

	return detail, nil
}

func (r Repository) Update(ctx context.Context, data tasks.Update) (entity.Tasks, error) {
	var detail entity.Tasks

	err := r.NewSelect().Model(&detail).Where("id = ?", data.Id).Scan(ctx)
	if err != nil {
		return entity.Tasks{}, err
	}

	if data.Status != nil {
		detail.Status = data.Status
	}
	if data.ProjectId != nil {
		detail.ProjectId = data.ProjectId
	}
	if data.Name != nil {
		detail.Name = data.Name
	}
	if data.Description != nil {
		detail.Description = data.Description
	}
	if data.AssignedTo != nil {
		detail.AssignedTo = data.AssignedTo
	}
	if data.Priority != nil {
		detail.Priority = data.Priority
	}
	if data.DueDate != nil {
		detail.DueDate = data.DueDate
	}

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)
	if err != nil {
		return entity.Tasks{}, err
	}

	return detail, nil
}

func (r Repository) Delete(ctx context.Context, data basic_service.Delete) error {
	return basic_repo.BasicDelete(ctx, data, &entity.Tasks{}, r.DB)
}
