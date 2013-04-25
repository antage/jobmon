package main

import (
	"jobmon/job"
	"jobmon/logger"
	"time"
)

type Server struct {
	logs *job.Logs
}

func newServer() *Server {
	server := new(Server)
	server.logs = job.NewLogs()
	return server
}

func (s *Server) close() {
	// do nothing
}

func (s *Server) StartNotification(args *job.StartNotification, logIdPtr *job.LogId) error {
	logger.Info("StartNotification: %s", *args)

	logId, entry := s.logs.NewLogEntry(args.JobId)
	*logIdPtr = logId

	entry.JobId = args.JobId
	entry.StartedAt = args.StartedAt
	entry.AliveAt = args.StartedAt

	return nil
}

func (s *Server) AliveNotification(args *job.AliveNotification, success *bool) error {
	logger.Info("AliveNotification: %s", *args)
	*success = false

	logEntryPtr, ok := s.logs.LogEntryById(args.LogId)
	if ok {
		logEntryPtr.AliveAt = time.Now()
		*success = true
	}
	return nil
}

func (s *Server) CompleteNotification(args *job.CompleteNotification, success *bool) error {
	logger.Info("CompleteNotification: %s", *args)
	*success = false
	logEntryPtr, ok := s.logs.LogEntryById(args.LogId)
	if ok {
		logEntryPtr.CompletedAt = args.CompletedAt
		logEntryPtr.Output = args.Output
		logEntryPtr.Success = args.Success
		*success = true

		if !logEntryPtr.Success {
			go sendNotification(s, logEntryPtr)
		}
	}
	return nil
}
