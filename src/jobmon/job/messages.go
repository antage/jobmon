package job

import (
	"time"
)

type StartNotification struct {
	JobId     *JobId
	StartedAt time.Time
}

type AliveNotification struct {
	LogId LogId
}

type CompleteNotification struct {
	LogId       LogId
	CompletedAt time.Time
	Output      []byte
	Success     bool
}
