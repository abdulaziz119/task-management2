package users

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uptrace/bun"
	"task-management/internal/entity"
	basic_repo "task-management/internal/repository/postgres/_basic_repo"
	basic_service "task-management/internal/service/_basic_service"
	"task-management/internal/service/users"
)

type Repository struct {
	*bun.DB
}

func (r Repository) GetAllUsers(ctx context.Context, filter users.Filter) ([]users.User, int, error) {
	query := `
        SELECT 
            id,
            full_name,
            email,
            role
        FROM users
        WHERE deleted_at IS NULL
        ORDER BY id`

	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}
	if filter.Offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *filter.Offset)
	}

	countQuery := "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL"
	var count int
	err := r.QueryRowContext(ctx, countQuery).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting users: %v", err)
	}

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var result []users.User
	for rows.Next() {
		var user users.User
		var id int64
		var fullName, email, role string
		err := rows.Scan(
			&id,
			&fullName,
			&email,
			&role,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning user row: %v", err)
		}
		user.Id = &id
		user.FullName = &fullName
		user.Email = &email
		user.Role = &role
		result = append(result, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating user rows: %v", err)
	}

	return result, count, nil
}

func (r Repository) GetTaskStats(ctx context.Context) (map[int64]users.TaskStats, error) {
	query := `
        SELECT 
            assigned_to,
            COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_tasks,
            COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_tasks,
            COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks
        FROM tasks 
        WHERE deleted_at IS NULL
        GROUP BY assigned_to`

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying task stats: %v", err)
	}
	defer rows.Close()

	stats := make(map[int64]users.TaskStats)
	for rows.Next() {
		var userId int64
		var pending, inProgress, completed int
		err := rows.Scan(
			&userId,
			&pending,
			&inProgress,
			&completed,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning task stats row: %v", err)
		}
		total := pending + inProgress + completed
		stats[userId] = users.TaskStats{
			PendingTasks:    &pending,
			InProgressTasks: &inProgress,
			CompletedTasks:  &completed,
			TaskCount:       &total,
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task stats rows: %v", err)
	}

	return stats, nil
}

func (r Repository) GetAll(ctx context.Context, filter users.Filter) ([]users.List, int, error) {
	usersList, count, err := r.GetAllUsers(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	taskStats, err := r.GetTaskStats(ctx)
	if err != nil {
		return nil, 0, err
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

func (r Repository) getUserDetails(ctx context.Context, userId int) (users.Detail, error) {
	query := `
		SELECT 
			id,
			full_name,
			email,
			role,
			created_at at time zone 'UTC'
		FROM users
		WHERE id = ? AND deleted_at IS NULL`

	var result users.Detail
	var id int64
	var fullName, email, role, createdAt string

	err := r.QueryRow(query, userId).Scan(
		&id,
		&fullName,
		&email,
		&role,
		&createdAt,
	)
	if err != nil {
		return users.Detail{}, fmt.Errorf("error getting user details: %w", err)
	}

	result.Id = &id
	result.FullName = &fullName
	result.Email = &email
	result.Role = &role
	result.CreatedAt = &createdAt

	return result, nil
}

func (r Repository) getTaskStatsCount(ctx context.Context, userId int) (users.TaskStats, error) {
	query := `
        SELECT 
            COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_tasks,
            COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_tasks,
            COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks,
            COUNT(*) as total_tasks
        FROM tasks
        WHERE deleted_at IS NULL AND assigned_to = ?
        GROUP BY assigned_to`

	var taskStats users.TaskStats
	var pending, inProgress, completed, total int

	err := r.QueryRow(query, userId).Scan(
		&pending,
		&inProgress,
		&completed,
		&total,
	)
	if err != nil {
		return users.TaskStats{}, fmt.Errorf("error getting task stats: %w", err)
	}

	taskStats.PendingTasks = &pending
	taskStats.InProgressTasks = &inProgress
	taskStats.CompletedTasks = &completed
	taskStats.TaskCount = &total

	return taskStats, nil
}

func (r Repository) getUserTasks(ctx context.Context, userId int) ([]users.TaskItem, error) {
	query := `
        SELECT 
            COALESCE(
                json_agg(
                    json_build_object(
                        'id', id,
                        'name', name,
                        'description', description,
                        'status', status,
                        'priority', priority,
                        'due_date', to_char(due_date at time zone 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
                        'created_at', to_char(created_at at time zone 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
                    ) ORDER BY created_at DESC
                ),
                '[]'::json
            ) as tasks
        FROM tasks
        WHERE deleted_at IS NULL AND assigned_to = ?`

	var tasksJson []byte
	err := r.QueryRow(query, userId).Scan(&tasksJson)
	if err != nil {
		return nil, fmt.Errorf("error getting user tasks: %w", err)
	}

	var tasks []users.TaskItem
	if err := json.Unmarshal(tasksJson, &tasks); err != nil {
		return nil, fmt.Errorf("error parsing tasks json: %w", err)
	}

	return tasks, nil
}

func (r Repository) GetById(ctx context.Context, userId int) (users.Detail, error) {
	result, err := r.getUserDetails(ctx, userId)
	if err != nil {
		return users.Detail{}, err
	}

	taskStats, err := r.getTaskStatsCount(ctx, userId)
	if err != nil {
		return users.Detail{}, err
	}

	tasks, err := r.getUserTasks(ctx, userId)
	if err != nil {
		return users.Detail{}, err
	}

	result.PendingTasks = taskStats.PendingTasks
	result.InProgressTasks = taskStats.InProgressTasks
	result.CompletedTasks = taskStats.CompletedTasks
	result.TaskCount = taskStats.TaskCount
	result.Tasks = &tasks

	return result, nil
}

func (r Repository) Create(ctx context.Context, data users.Create) (entity.User, error) {
	var detail entity.User

	detail.Email = data.Email
	detail.Password = data.Password
	detail.FullName = data.FullName
	detail.Role = data.Role

	_, err := r.NewInsert().Model(&detail).Exec(ctx)

	return detail, err
}

func (r Repository) Update(ctx context.Context, data users.Update) (entity.User, error) {
	var detail entity.User

	err := r.NewSelect().Model(&detail).Where("id = ?", data.Id).Scan(ctx)
	if err != nil {
		return entity.User{}, err
	}

	if data.FullName != nil {
		detail.FullName = data.FullName
	}
	if data.Role != nil {
		detail.Role = data.Role
	}
	if data.Email != nil {
		detail.Email = data.Email
	}
	if data.Password != nil {
		detail.Password = data.Password
	}

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	return detail, err
}

func (r Repository) Delete(ctx context.Context, data basic_service.Delete) error {
	return basic_repo.BasicDelete(ctx, data, &entity.User{}, r.DB)
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}
