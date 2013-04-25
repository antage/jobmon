package main

import (
	"flag"
	"fmt"
	"jobmon/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configFilename = flag.String("c", fmt.Sprintf("%s/%s", CONF_DIR, "jobmond.json"), "Configuration file location")

var serverInstance *Server

func handleKill() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		for s := range c {
			if serverInstance != nil {
				logger.Info("closing server")
				serverInstance.close()
			}
			logger.Info("exited, got %s signal", s)
			exit(0)
		}
	}()
}

func exit(errcode int) {
	logger.Close()
	os.Exit(errcode)
}

func main() {
	flag.Parse()
	logger.Init(fmt.Sprintf("jobmond[%d]", os.Getpid()))
	logger.Info("started")
	handleKill()

	err := parseConfigFile(*configFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't load configuration file: %s\n", err.Error())
		exit(1)
	}

	serverInstance = newServer()

	err = runRpcServer(serverInstance)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't run RPC server: %s", err.Error())
		exit(1)
	}

	httpServer := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Web.Listen, config.Web.Port),
		Handler:        routes(serverInstance),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 65536,
	}
	go httpServer.ListenAndServe()

	for {
		time.Sleep(5 * time.Second)
	}
}
