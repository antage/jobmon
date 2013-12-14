package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jobmon/job"
	"jobmon/logger"
	"net/http"
	"strconv"
)

type jsonLogEntry struct {
	Id      job.LogId
	JobId   *job.JobId
	Output  string
	Success bool
}

func logEntryToJson(l *job.LogEntry) *jsonLogEntry {
	return &jsonLogEntry{
		Id:      l.Id,
		JobId:   l.JobId,
		Output:  string(l.Output),
		Success: l.Success,
	}
}

func logEntryJson(server *Server, resp http.ResponseWriter, req *http.Request) {
	logger.Info("process http request '%s' from %s", req.RequestURI, req.RemoteAddr)
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")

	vars := mux.Vars(req)
	idStr, ok := vars["id"]
	if !ok {
		logger.Error("can't get log entry ID from GET request")
		http.Error(resp, "{ error: \"need ID\" }", http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("can't convert log entry ID to integer number: %s", err.Error())
		http.Error(resp, fmt.Sprintf("{ error: \"%s\" }", err.Error()), http.StatusInternalServerError)
		return
	}

	logEntry, ok := server.logs.LogEntryById(job.LogId(id))
	if !ok {
		logger.Error("can't get log entry with ID = %d", id)
		http.Error(resp, "{ error: \"not found\" }", http.StatusNotFound)
		return
	}

	output, err := json.Marshal(logEntryToJson(logEntry))
	if err != nil {
		logger.Error("can't encode log entry (id=%d) in json: %s", logEntry.Id, err.Error())
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
