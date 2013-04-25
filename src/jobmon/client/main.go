package main

import (
	"bytes"
	"flag"
	"fmt"
	"jobmon/job"
	"jobmon/logger"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"os/user"
	"time"
)

const CONN_TIMEOUT = 60 * time.Second

var shell = flag.String("s", "/bin/sh", "shell path")
var configFilename = flag.String("c", fmt.Sprintf("%s/%s", CONF_DIR, "jobmon.json"), "configuration file location")

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

	err = parseConfigFile(*configFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't load configuration file: %s\n", err.Error())
		exit(4)
	}

	notifyJobmond := false
	var jobmond *rpc.Client

	caddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", config.RPC.Host, config.RPC.Port))
	if err != nil {
		logger.Error("can't resolve RPC server address: %s", err.Error())
		fmt.Fprintf(os.Stderr, "can't resolve RPC server address: %s\n", err.Error())
	} else {
		conn, err := net.DialTCP("tcp4", nil, caddr)
		if err != nil {
			logger.Error("can't connect to RPC server: %s", err.Error())
			fmt.Fprintf(os.Stderr, "can't connect to RPC server: %s\n", err.Error())
		} else {
			conn.SetDeadline(time.Now().Add(CONN_TIMEOUT))

			jobmond = rpc.NewClient(conn)
			notifyJobmond = true
		}
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
	if notifyJobmond {
		arg := &job.StartNotification{
			JobId:     &jobId,
			StartedAt: time.Now()}
		err = jobmond.Call("Server.StartNotification", arg, &logId)
		if err != nil {
			logger.Error("can't send StartNotification message: %s", err.Error())
		}
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
			if notifyJobmond {
				arg := &job.AliveNotification{LogId: logId}
				var success bool
				err = jobmond.Call("Server.AliveNotification", arg, &success)
				if err != nil {
					logger.Error("can't send AliveNotification message: %s", err.Error())
				}
				if !success {
					logger.Error("server can't handle AliveNotification: success is false")
				}
			}
		}
	}

	// notify jobmond the command complete and send exit code, stdout and stderr logs
	if notifyJobmond {
		arg := &job.CompleteNotification{
			LogId:       logId,
			CompletedAt: time.Now(),
			Output:      output.Bytes(),
			Success:     cmd_err == nil}
		var success bool
		err = jobmond.Call("Server.CompleteNotification", arg, &success)
		if err != nil {
			logger.Error("can't send CompleteNotification message: %s", err.Error())
		}
		if !success {
			logger.Error("server can't handle CompleteNotification: success is false")
		}
	}
}
