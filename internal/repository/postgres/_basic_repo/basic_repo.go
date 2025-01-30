package basic_repo

import (
	"context"
	"github.com/uptrace/bun"
	"time"
)

func BasicDelete(ctx context.Context, data Delete, table interface{}, r *bun.DB) error {
	_, err := r.NewUpdate().
		Model(table).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", *data.Id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return err
}
