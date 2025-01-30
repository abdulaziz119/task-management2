package entity

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	basicEntity
	FullName *string `bun:"full_name,notnull"`
	Email    *string `bun:"email,notnull"`
	Role     *string `bun:"role,notnull"`
	Password *string `bun:"password,notnull"`
}
