package projects

import (
	"context"
	"database/sql"
	"fmt"
	"task-management/internal/entity"
	basic_service "task-management/internal/service/_basic_service"
	"task-management/internal/service/projects"
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func (r Repository) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return r.DB.QueryRowContext(ctx, query, args...)
}

func (r Repository) buildTaskStatsQuery() string {
	return `
		WITH task_stats AS (
			SELECT 
				project_id,
				COUNT(*) as total_tasks,
				COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks,
				COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_tasks
			FROM tasks
			WHERE deleted_at IS NULL
			GROUP BY project_id
		)
	`
}

func (r Repository) buildProjectsBaseQuery(whereClause string) string {
	return fmt.Sprintf(`
		SELECT 
			p.*,
			COALESCE(ts.total_tasks, 0) as total_tasks,
			COALESCE(
				CASE 
					WHEN ts.total_tasks > 0 THEN
						(
							(COALESCE(ts.completed_tasks, 0)::numeric * 100 + 
							COALESCE(ts.in_progress_tasks, 0)::numeric * 50) / 
							(ts.total_tasks::numeric * 100) * 100
						)::numeric(10,1)
					ELSE 0
				END
			, 0) as progress
		FROM projects p
		LEFT JOIN task_stats ts ON p.id = ts.project_id
		WHERE p.deleted_at IS NULL
		%s
	`, whereClause)
}

func (r Repository) buildFinalSelectQuery() string {
	return `
		SELECT 
			id,
			name,
			description,
			owner_id,
			total_tasks,
			progress
		FROM projects_with_stats
	`
}

func (r Repository) buildWhereAndParams(filter projects.Filter) (string, []interface{}) {
	var whereClause string
	var params []interface{}

	if filter.OwnerId != nil {
		whereClause = "AND p.owner_id = ?"
		params = append(params, *filter.OwnerId)
	}

	return whereClause, params
}

func (r Repository) buildLimitOffset(filter projects.Filter, params []interface{}) (string, []interface{}) {
	var limitOffsetClause string

	if filter.Limit != nil {
		limitOffsetClause = "LIMIT ?"
		params = append(params, *filter.Limit)

		if filter.Offset != nil {
			limitOffsetClause += " OFFSET ?"
			params = append(params, *filter.Offset)
		}
	}

	return limitOffsetClause, params
}

func (r Repository) buildProjectsQuery(filter projects.Filter) (string, []interface{}) {
	whereClause, params := r.buildWhereAndParams(filter)
	limitOffsetClause, params := r.buildLimitOffset(filter, params)

	query := fmt.Sprintf(`
		%s,
		projects_with_stats AS (
			%s
		)
		%s
		%s
	`,
		r.buildTaskStatsQuery(),
		r.buildProjectsBaseQuery(whereClause),
		r.buildFinalSelectQuery(),
		limitOffsetClause,
	)

	return query, params
}

func (r Repository) GetProjectsWithStats(ctx context.Context, filter projects.Filter) ([]projects.List, error) {
	var result []projects.List

	query, params := r.buildProjectsQuery(filter)
	rows, err := r.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item projects.List
		var name, description *string
		var ownerId *int
		var totalTasks int
		var progress float64

		err = rows.Scan(
			&item.Id,
			&name,
			&description,
			&ownerId,
			&totalTasks,
			&progress,
		)
		if err != nil {
			return nil, err
		}

		if name != nil {
			item.Name = *name
		}
		if description != nil {
			item.Description = *description
		}
		if ownerId != nil {
			item.OwnerId = *ownerId
		}
		item.TotalTasks = totalTasks
		item.Progress = progress

		result = append(result, item)
	}

	return result, nil
}

func (r Repository) GetProjectsCount(ctx context.Context, filter projects.Filter) (int, error) {
	var count int

	query := `
		SELECT COUNT(*) 
		FROM projects p
		WHERE p.deleted_at IS NULL
	`

	var params []interface{}

	if filter.OwnerId != nil {
		query += " AND p.owner_id = ?"
		params = append(params, *filter.OwnerId)
	}

	err := r.DB.QueryRowContext(ctx, query, params...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r Repository) buildTaskStatsForFindOneQuery() string {
	return `
		WITH task_stats AS (
			SELECT 
				project_id,
				COUNT(*) as total_tasks,
				COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks,
				COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_tasks
			FROM tasks
			WHERE deleted_at IS NULL
			GROUP BY project_id
		)
	`
}

func (r Repository) buildProjectWithStatsFindOneQuery() string {
	return `
		SELECT 
			p.id,
			p.name,
			p.description,
			p.owner_id,
			COALESCE(ts.total_tasks, 0) as total_tasks,
			COALESCE(
				CASE 
					WHEN ts.total_tasks > 0 THEN
						(
							(COALESCE(ts.completed_tasks, 0)::numeric * 100 + 
							COALESCE(ts.in_progress_tasks, 0)::numeric * 50) / 
							(ts.total_tasks::numeric * 100) * 100
						)::numeric(10,1)
					ELSE 0
				END
			, 0) as progress
		FROM projects p
		LEFT JOIN task_stats ts ON p.id = ts.project_id
		WHERE p.id = ? AND p.deleted_at IS NULL
	`
}

func (r Repository) buildFindOneQuery() string {
	return fmt.Sprintf(`
		%s
		%s
	`,
		r.buildTaskStatsForFindOneQuery(),
		r.buildProjectWithStatsFindOneQuery(),
	)
}

func (r Repository) scanProjectDetailForFindOne(row *sql.Row) (projects.Detail, error) {
	var detail projects.Detail
	var name, description *string
	var ownerId *int
	var totalTasks int
	var progress float64

	err := row.Scan(
		&detail.Id,
		&name,
		&description,
		&ownerId,
		&totalTasks,
		&progress,
	)
	if err != nil {
		return projects.Detail{}, err
	}

	if name != nil {
		detail.Name = *name
	}
	if description != nil {
		detail.Description = *description
	}
	if ownerId != nil {
		detail.Owner_id = *ownerId
	}

	detail.TaskStats = projects.TaskStats{
		TotalTasks: totalTasks,
		Progress:   progress,
	}

	return detail, nil
}

func (r Repository) GetById(ctx context.Context, id int) (projects.Detail, error) {
	query := r.buildFindOneQuery()
	row := r.QueryRowContext(ctx, query, id)
	return r.scanProjectDetailForFindOne(row)
}

func (r Repository) Create(ctx context.Context, data projects.Create) (entity.Projects, error) {
	var project entity.Projects

	query := `
		INSERT INTO projects (name, description, owner_id, created_at)
		VALUES (?, ?, ?, ?)
		RETURNING id, name, description, owner_id, created_at
	`

	now := time.Now()
	err := r.DB.QueryRowContext(ctx, query,
		data.Name,
		data.Description,
		data.Owner_id,
		now,
	).Scan(
		&project.Id,
		&project.Name,
		&project.Description,
		&project.OwnerId,
		&project.CreatedAt,
	)

	if err != nil {
		return entity.Projects{}, err
	}

	return project, nil
}

func (r Repository) Update(ctx context.Context, data projects.Update) (entity.Projects, error) {
	var project entity.Projects

	query := `
		UPDATE projects 
		SET 
			name = COALESCE(?, name),
			description = COALESCE(?, description),
			owner_id = COALESCE(?, owner_id)
		WHERE id = ? AND deleted_at IS NULL
		RETURNING id, name, description, owner_id, created_at
	`

	err := r.DB.QueryRowContext(ctx, query,
		data.Name,
		data.Description,
		data.Owner_id,
		data.Id,
	).Scan(
		&project.Id,
		&project.Name,
		&project.Description,
		&project.OwnerId,
		&project.CreatedAt,
	)

	if err != nil {
		return entity.Projects{}, err
	}

	return project, nil
}

func (r Repository) Delete(ctx context.Context, data basic_service.Delete) error {
	query := `
		UPDATE projects 
		SET deleted_at = ? 
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.DB.ExecContext(ctx, query, time.Now(), data.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB}
}
