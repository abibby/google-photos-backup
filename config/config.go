package config

import (
	"context"

	"github.com/abibby/salusa/env"
	"github.com/joho/godotenv"
)

var Port int
var DBPath string

var GoogleClientID string
var GoogleClientSecret string

var PhotoDir string

func Load(ctx context.Context) error {
	err := godotenv.Load("./.env")
	if err != nil {
		return err
	}

	Port = env.Int("PORT", 6900)
	DBPath = env.String("DATABASE_PATH", "./db.sqlite")
	GoogleClientID = env.String("CLIENT_ID", "")
	GoogleClientSecret = env.String("CLIENT_SECRET", "")
	PhotoDir = env.String("PHOTO_DIR", "")

	return nil
}
