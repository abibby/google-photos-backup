package events

import (
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/event/cron"
)

type BackupEvent struct {
	cron.BaseEvent
}

var _ event.Event = (*BackupEvent)(nil)

func (e *BackupEvent) Type() event.EventType {
	return "gpb:backup"
}

// func init() {
// 	kernel.RegisterEvent(&LogEvent{})
// }
