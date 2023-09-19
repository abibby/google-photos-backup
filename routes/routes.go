package routes

import (
	"github.com/abibby/google-photos-backup/app/handlers"
	"github.com/abibby/google-photos-backup/resources"
	"github.com/abibby/salusa/fileserver"
	"github.com/abibby/salusa/router"
)

func InitRoutes(r *router.Router) {
	r.Get("/backup", handlers.Backup)

	r.Handle("/", fileserver.WithFallback(resources.Content, "dist", "index.html", nil))
}
