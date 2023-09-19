package server

import (
	"fmt"
	"net/http"

	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/google-photos-backup/config"
	"github.com/abibby/google-photos-backup/database"
	"github.com/abibby/google-photos-backup/routes"
)

func Serve() error {
	r := router.New()

	r.Use(request.WithDB(database.DB))

	routes.InitRoutes(r)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
}
