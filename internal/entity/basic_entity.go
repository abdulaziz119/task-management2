package entity

import "time"

type basicEntity struct {
	Id        int        `json:"id" bun:"id,pk,autoincrement"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	UpdateAt  *time.Time `json:"updated_at" bun:"updated_at"`
}
