package routes

import (
	"github.com/abibby/google-photos-backup/app/handlers"
	"github.com/abibby/google-photos-backup/config"
	"github.com/abibby/google-photos-backup/database"
	"github.com/abibby/google-photos-backup/resources"
	"github.com/abibby/salusa/fileserver"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
)

func InitRoutes(r *router.Router) {
	r.Use(request.HandleErrors())
	r.Use(request.WithDB(database.DB))

	r.Get("/backup", handlers.Backup)
	r.Get("/gauth", handlers.Login)

	r.Handle("/", fileserver.WithFallback(resources.Content, "dist", "index.html", map[string]string{
		"client_id": config.GoogleClientID,
	}))
}
