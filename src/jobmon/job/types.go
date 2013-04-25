package job

import (
	"fmt"
	"time"
)

type JobId struct {
	Hostname string
	Username string
	Name     string
}

func (jobId JobId) String() string {
	return fmt.Sprintf("{Hostname: %s, Username: %s, Jobname: %s}", jobId.Hostname, jobId.Username, jobId.Name)
}

type LogId int64

type LogEntry struct {
	Id          LogId
	JobId       *JobId
	StartedAt   time.Time
	AliveAt     time.Time
	CompletedAt time.Time
	Output      []byte
	Success     bool
}
