package job

import (
	"sync"
)

type Logs struct {
	sync.RWMutex

	seq    seqInt64
	logIds map[JobId][]LogId
	logs   map[LogId]*LogEntry
}

func NewLogs() (l *Logs) {
	l = new(Logs)
	l.logIds = make(map[JobId][]LogId)
	l.logs = make(map[LogId]*LogEntry)
	return
}

func (l *Logs) NewLogEntry(jobIdPtr *JobId) (LogId, *LogEntry) {
	l.Lock()
	defer l.Unlock()

	jobId := *jobIdPtr
	logId := LogId(l.seq.next())

	logEntryPtr := new(LogEntry)
	logEntryPtr.Id = logId
	logEntryPtr.JobId = jobIdPtr

	l.logs[logId] = logEntryPtr

	logIds, ok := l.logIds[jobId]
	if !ok {
		logIds = make([]LogId, 0, 100)
	}

	if len(logIds) > 99 {
		copy(logIds, logIds[(len(logIds)-99):])
		logIds = logIds[0:99]
	}

	l.logIds[jobId] = append(logIds, logId)

	return logId, logEntryPtr
}

func (l *Logs) LogEntryById(logId LogId) (*LogEntry, bool) {
	l.RLock()
	defer l.RUnlock()

	logEntry, ok := l.logs[logId]
	return logEntry, ok
}

func (l *Logs) JobIds() []JobId {
	l.RLock()
	defer l.RUnlock()

	jobIds := make([]JobId, 0, len(l.logIds))
	for jobId, _ := range l.logIds {
		jobIds = append(jobIds, jobId)
	}
	return jobIds
}

func (l *Logs) LastLogEntryByJobId(jobIdPtr *JobId) *LogEntry {
	l.RLock()
	defer l.RUnlock()

	logIds, ok := l.logIds[*jobIdPtr]
	if ok {
		logId := logIds[len(logIds)-1]
		logEntry, ok2 := l.LogEntryById(logId)
		if ok2 {
			return logEntry
		} else {
			return nil
		}
	} else {
		return nil
	}
	panic("unreachable")
}
