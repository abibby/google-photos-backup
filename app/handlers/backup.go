package handlers

import (
	"github.com/abibby/google-photos-backup/app/events"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
)

type BackupRequest struct {
}
type BackupResponse struct {
}

var Backup = request.Handler(func(r *BackupRequest) (*BackupResponse, error) {
	kernel.Dispatch(&events.BackupEvent{})
	return &BackupResponse{}, nil
})
