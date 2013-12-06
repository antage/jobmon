package main

import (
	"bytes"
	"flag"
	"fmt"
	"jobmon/client/config"
	"jobmon/job"
	"jobmon/logger"
	"os"
	"os/exec"
	"os/user"
	"time"
)

var shell = flag.String("s", "/bin/sh", "shell path")
var configFilename = flag.String("c", fmt.Sprintf("%s/%s", config.CONF_DIR, "jobmon.json"), "configuration file location")

var configuration *config.Main

func exit(errcode int) {
	logger.Close()
	os.Exit(errcode)
}

func main() {
	flag.Parse()
	logger.Init(fmt.Sprintf("jobmon[%d]", os.Getpid()))
	defer logger.Close()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "You should give job name and command\n")
		exit(2)
	}

	var jobId job.JobId
	jobId.Name = flag.Arg(0)

	command := flag.Arg(1)

	user, err := user.Current()
	if err != nil {
		logger.Error("can't get username: ", err.Error())
		fmt.Fprintf(os.Stderr, "can't get username: %s\n", err.Error())
		exit(4)
	}
	jobId.Username = user.Username

	jobId.Hostname, err = os.Hostname()
	if err != nil {
		logger.Error("can't get hostname: ", err.Error())
		fmt.Fprintf(os.Stderr, "can't get hostname: %s\n", err.Error())
		exit(4)
	}

	configuration, err = config.ParseFile(*configFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't load configuration file: %s\n", err.Error())
		exit(4)
	}

	logger.Info("started with JobId: %s, command: (%s)", jobId, command)
	cmd_exec := exec.Command(*shell, []string{"-c", command}...)

	var output bytes.Buffer
	cmd_exec.Stderr = &output
	cmd_exec.Stdout = &output

	err = cmd_exec.Start()
	if err != nil {
		logger.Error("can't start command: %s", err.Error())
		fmt.Fprintf(os.Stderr, "can't start command: %s\n", err.Error())
		exit(1)
	}

	// notify jobmond the command started
	var logId job.LogId
	var serverAlive bool

	logId, err = rpcStartNotification(jobId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't send StartNotification to RPC server: %s\n", err.Error())
		serverAlive = false
	} else {
		serverAlive = true
	}

	cmd_complete_ch := make(chan error)
	go func() {
		cmd_complete_ch <- cmd_exec.Wait()
	}()

	var cmd_err error
command_waiting:
	for {
		select {
		case cmd_err = <-cmd_complete_ch:
			logger.Info("completed")
			break command_waiting
		case <-time.After(1 * time.Minute):
			// notify jobmond the command is alive
			if serverAlive {
				rpcAliveNotification(logId)
			}
		}
	}

	// notify jobmond the command complete and send exit code, stdout and stderr logs
	if serverAlive {
		rpcCompleteNotification(logId, time.Now(), output.Bytes(), cmd_err == nil)
	}
}
