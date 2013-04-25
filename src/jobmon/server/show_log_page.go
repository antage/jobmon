package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"jobmon/job"
	"jobmon/logger"
	"net/http"
	"strconv"
)

var showLogHtmlTemplate *template.Template

func init() {
	t := template.New("show_log.html.template")
	showLogHtmlTemplate = template.Must(t.ParseFiles(fmt.Sprintf("%s/%s", VIEWS_DIR, "show_log.html.template")))
}

func showLogPage(server *Server, resp http.ResponseWriter, req *http.Request) {
	logger.Info("process http request '%s' from %s", req.RequestURI, req.RemoteAddr)
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")

	vars := mux.Vars(req)
	idStr, ok := vars["id"]
	logger.Info("idStr: %s", idStr)
	if !ok {
		http.Error(resp, "404 page not found", http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(idStr)
	logger.Info("id: %d", id)
	if err != nil {
		http.Error(resp, fmt.Sprintf("500 error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	logEntry, ok := server.logs.LogEntryById(job.LogId(id))
	if !ok {
		http.Error(resp, "404 page not found", http.StatusNotFound)
		return
	}

	data := struct {
		Id     string
		Output string
	}{
		Id:     fmt.Sprintf("%d", logEntry.Id),
		Output: string(logEntry.Output),
	}

	err = showLogHtmlTemplate.Execute(resp, data)
	if err != nil {
		logger.Error("can't render template: %s", err)
	}
}
