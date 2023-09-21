package models

import (
	"context"
	"time"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/jmoiron/sqlx"
)

//go:generate spice generate:migration
type Photo struct {
	model.BaseModel

	ID    int    `json:"id"    db:"id,primary,autoincrement"`
	Email string `json:"email" db:"email"`

	AccessToken  string    `json:"-" db:"access_token"`
	ExpiresAt    time.Time `json:"-" db:"expires_in"`
	RefreshToken string    `json:"-" db:"refresh_token"`
}

func PhotoQuery(ctx context.Context) *builder.Builder[*Photo] {
	return builder.From[*Photo]().WithContext(ctx)
}

func (u *Photo) Save(tx *sqlx.Tx) error {
	return model.Save(tx, u)
}
