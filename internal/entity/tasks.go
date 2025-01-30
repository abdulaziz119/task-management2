package entity

import (
	"github.com/uptrace/bun"
)

type Tasks struct {
	bun.BaseModel `bun:"table:tasks"`

	basicEntity
	ProjectId   *int    `json:"project_id" bun:"project_id"`
	Name        *string `json:"name" bun:"name"`
	Description *string `json:"description" bun:"description"`
	AssignedTo  *int    `json:"assigned_to" bun:"assigned_to"`
	Status      *string `json:"status" bun:"status"`
	Priority    *string `json:"priority" bun:"priority"`
	DueDate     *string `json:"due_date" bun:"due_date"`
}
