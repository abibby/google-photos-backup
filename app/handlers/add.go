package handlers

import (
	"github.com/abibby/google-photos-backup/app/events"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
)

type backupRequest struct {
}
type backupResponse struct {
}

var Backup = request.Handler(func(r *backupRequest) (*backupResponse, error) {
	kernel.Dispatch(&events.BackupEvent{})
	return &backupResponse{}, nil
})
