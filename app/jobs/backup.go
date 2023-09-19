package jobs

import (
	"log"

	"github.com/abibby/google-photos-backup/app/events"
)

func BackupJob(e *events.BackupEvent) error {
	log.Print("backup")
	return nil
}
