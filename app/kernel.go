package app

import (
	"github.com/abibby/google-photos-backup/app/jobs"
	"github.com/abibby/google-photos-backup/config"
	"github.com/abibby/google-photos-backup/database"
	"github.com/abibby/google-photos-backup/routes"
	"github.com/abibby/salusa/kernel"
)

var Kernel = kernel.NewDefaultKernel(
	kernel.Port(func() int {
		return config.Port
	}),
	kernel.Bootstrap(
		config.Load,
		database.Init,
	),
	kernel.Services(
	// cron.Service().
	// 	Schedule("* * * * *", &events.BackupEvent{}),
	),
	kernel.Listeners(
		kernel.NewListener(jobs.BackupJob),
	),
	kernel.InitRoutes(routes.InitRoutes),
)
