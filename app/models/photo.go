package models

import (
	"context"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/jmoiron/sqlx"
)

//go:generate spice generate:migration
type Photo struct {
	model.BaseModel

	ID      int    `json:"id"         db:"id,primary,autoincrement"`
	UserID  int    `json:"user_id"    db:"user_id"`
	PhotoID string `json:"photo_id"   db:"photo_id,unique"`
}

func PhotoQuery(ctx context.Context) *builder.Builder[*Photo] {
	return builder.From[*Photo]().WithContext(ctx)
}

func (u *Photo) Save(tx *sqlx.Tx) error {
	return model.Save(tx, u)
}
