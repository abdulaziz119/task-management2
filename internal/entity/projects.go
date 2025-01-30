package entity

import "github.com/uptrace/bun"

type Projects struct {
	bun.BaseModel `bun:"table:projects"`

	basicEntity
	Name        *string `json:"name" bun:"name"`
	Description *string `json:"description" bun:"description"`
	OwnerId     *int    `json:"owner_id" bun:"owner_id"`
}
