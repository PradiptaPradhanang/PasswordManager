package backUpDB

import (
	"passmana/model"
	"sync"
)

type BackUpEvent struct {
	Action string
	Data   model.Cred
}
type BackupManager struct {
	mu     sync.Mutex
	events chan BackUpEvent
	stop   chan struct{}
}

func (bm *BackupManager) Send(evt BackUpEvent) {
	bm.events <- evt
}

func (bm *BackupManager) Stop() {
	close(bm.stop)
}
