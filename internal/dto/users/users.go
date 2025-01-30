package users

import "time"

type Filter struct {
	Limit  *int
	Offset *int
}

type Create struct {
	FullName *string `json:"full_name" bun:"full_name"`
	Email    *string `json:"email" bun:"username,unique,notnull"`
	Role     *string `json:"role" validate:"required,oneof=manager worker"`
	Password *string `json:"password" bun:"password"`
}

type Update struct {
	Id       *int    `json:"id" form:"id"`
	FullName *string `json:"full_name" bun:"full_name"`
	Email    *string `json:"email" bun:"username,unique,notnull"`
	Role     *string `json:"role" bun:"role"`
	Password *string `json:"password" bun:"password"`
}

type User struct {
	Id       *int64  `json:"id"`
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
	Role     *string `json:"role"`
}

type TaskStats struct {
	PendingTasks    *int `json:"pending_tasks"`
	InProgressTasks *int `json:"in_progress_tasks"`
	CompletedTasks  *int `json:"completed_tasks"`
	TaskCount       *int `json:"task_count"`
}

type List struct {
	Id              *int64  `json:"id"`
	FullName        *string `json:"full_name"`
	Email           *string `json:"email"`
	Role            *string `json:"role"`
	PendingTasks    *int    `json:"pending_tasks"`
	InProgressTasks *int    `json:"in_progress_tasks"`
	CompletedTasks  *int    `json:"completed_tasks"`
	TaskCount       *int    `json:"task_count"`
}

type TaskItem struct {
	Id          *int       `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Detail struct {
	Id              *int64      `json:"id"`
	FullName        *string     `json:"full_name"`
	Email           *string     `json:"email"`
	Role            *string     `json:"role"`
	PendingTasks    *int        `json:"pending_tasks"`
	InProgressTasks *int        `json:"in_progress_tasks"`
	CompletedTasks  *int        `json:"completed_tasks"`
	TaskCount       *int        `json:"task_count"`
	CreatedAt       *string     `json:"created_at"`
	UpdatedAt       *string     `json:"updated_at"`
	Tasks           *[]TaskItem `json:"tasks"`
}
