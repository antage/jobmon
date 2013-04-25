package main

import (
	"fmt"
	"html/template"
	"jobmon/job"
	"jobmon/logger"
	"net/http"
)

type HumanLogEntry struct {
	Id         string
	Success    bool
	StartedAt  string
	Duration   string
	Processing bool
}

var indexHtmlTemplate *template.Template

func init() {
	t := template.New("index.html.template")
	indexHtmlTemplate = template.Must(t.ParseFiles(fmt.Sprintf("%s/%s", VIEWS_DIR, "index.html.template")))
}

func humanizeLogEntry(logEntry *job.LogEntry) *HumanLogEntry {
	h := new(HumanLogEntry)
	h.Id = fmt.Sprintf("%d", logEntry.Id)
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

func indexPage(server *Server, resp http.ResponseWriter, req *http.Request) {
	logger.Info("process http request '%s' from %s", req.RequestURI, req.RemoteAddr)
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")

	jobs := make(map[job.JobId]*HumanLogEntry)
	for _, jobId := range server.logs.JobIds() {
		jobs[jobId] = humanizeLogEntry(server.logs.LastLogEntryByJobId(&jobId))
	}

	data := struct {
		Jobs map[job.JobId]*HumanLogEntry
	}{
		Jobs: jobs,
	}

	err := indexHtmlTemplate.Execute(resp, data)
	if err != nil {
		logger.Error("can't render template: %s", err)
	}
}
