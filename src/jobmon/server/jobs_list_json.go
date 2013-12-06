package main

import (
	"encoding/json"
	"fmt"
	"jobmon/job"
	"jobmon/logger"
	"net/http"
)

type HumanLogEntry struct {
	Id         string
	JobId      *job.JobId
	Success    bool
	StartedAt  string
	Duration   string
	Processing bool
}

func humanizeLogEntry(logEntry *job.LogEntry) *HumanLogEntry {
	h := new(HumanLogEntry)
	h.Id = fmt.Sprintf("%d", logEntry.Id)
	h.JobId = logEntry.JobId
	h.StartedAt = logEntry.StartedAt.Format("2006-01-02 15:04:05")
	if logEntry.CompletedAt.IsZero() {
		h.Processing = true
		h.Duration = logEntry.AliveAt.Sub(logEntry.StartedAt).String()
		h.Success = true
	} else {
		h.Processing = false
		h.Duration = logEntry.CompletedAt.Sub(logEntry.StartedAt).String()
		h.Success = logEntry.Success
	}
	return h
}

func jobsListJson(server *Server, resp http.ResponseWriter, req *http.Request) {
	logger.Info("process http request '%s' from %s", req.RequestURI, req.RemoteAddr)
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")

	jobIds := server.logs.JobIds()
	jobs := make([]*HumanLogEntry, 0, len(jobIds))
	for _, jobId := range jobIds {
		jobs = append(jobs, humanizeLogEntry(server.logs.LastLogEntryByJobId(&jobId)))
	}

	output, err := json.Marshal(jobs)
	if err != nil {
		logger.Error("can't encode jobs list in json: %s", err.Error())
		http.Error(resp, fmt.Sprintf("{ error: \"%s\" }", err.Error()), http.StatusInternalServerError)
		return
	}
	_, err = resp.Write(output)
	if err != nil {
		logger.Error("can't write output: %s", err.Error())
		http.Error(resp, fmt.Sprintf("{ error: \"%s\" }", err.Error()), http.StatusInternalServerError)
		return
	}

}
