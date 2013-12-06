package main

import (
	"fmt"
	"jobmon/logger"
	"net/http"
)

func appPage(resp http.ResponseWriter, req *http.Request) {
	logger.Info("process http request '%s' from %s", req.RequestURI, req.RemoteAddr)
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(resp, req, fmt.Sprintf("%s/app.html", ASSETS_DIR))
}
