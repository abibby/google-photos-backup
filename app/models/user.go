package models

import (
	"context"
	"time"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/jmoiron/sqlx"
)

//go:generate spice generate:migration
type User struct {
	model.BaseModel

	ID    int    `json:"id"    db:"id,primary,autoincrement"`
	Email string `json:"email" db:"email"`

	AccessToken  string    `json:"-" db:"access_token"`
	ExpiresAt    time.Time `json:"-" db:"expires_in"`
	RefreshToken string    `json:"-" db:"refresh_token"`

	FinishedInitialFetch bool `json:"finished_initial_fetch" db:"finished_initial_fetch"`
}

func UserQuery(ctx context.Context) *builder.Builder[*User] {
	return builder.From[*User]().WithContext(ctx)
}

func (u *User) Save(tx *sqlx.Tx) error {
	return model.Save(tx, u)
}
