package main

import (
	"fmt"
	"jobmon/client/config"
	"jobmon/job"
	"jobmon/logger"
	"net"
	"net/rpc"
	"time"
)

const CONN_TIMEOUT = 60 * time.Second

type rpcClient struct {
	*rpc.Client
	conn net.Conn
}

func newFromConfig(c *config.Main) (*rpcClient, error) {
	conn, err := net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", c.RPC.Host, c.RPC.Port), CONN_TIMEOUT)
	if err != nil {
		logger.Error("can't resolve connect to server: %s", err.Error())
		return nil, err
	}

	return &rpcClient{rpc.NewClient(conn), conn}, nil
}

func (cl *rpcClient) close() {
	err := cl.Close()
	if err != nil {
		logger.Error("can't close RPC client: %s", err.Error())
		return
	}
	err = cl.conn.Close()
	if err != nil {
		logger.Error("can't close connection: %s", err.Error())
		return
	}
	return
}

func rpcStartNotification(jobId job.JobId) (logId job.LogId, err error) {
	cl, err := newFromConfig(configuration)
	if err != nil {
		logger.Error("can't create RPC client: %s", err.Error())
		return job.LogId(0), err
	}
	defer cl.close()

	arg := &job.StartNotification{
		JobId:     &jobId,
		StartedAt: time.Now(),
	}

	cl.conn.SetDeadline(time.Now().Add(CONN_TIMEOUT))
	err = cl.Call("Server.StartNotification", arg, &logId)
	if err != nil {
		return job.LogId(0), err
	}

	return logId, nil
}

func rpcAliveNotification(logId job.LogId) {
	cl, err := newFromConfig(configuration)
	if err != nil {
		logger.Error("can't create RPC client: %s", err.Error())
		return
	}
	defer cl.close()

	arg := &job.AliveNotification{LogId: logId}

	var success bool

	cl.conn.SetDeadline(time.Now().Add(CONN_TIMEOUT))
	err = cl.Call("Server.AliveNotification", arg, &success)

	if err != nil {
		logger.Error("can't send AliveNotification message: %s", err.Error())
		return
	}

	if !success {
		logger.Error("server can't handle AliveNotification: success is false")
	}
}

func rpcCompleteNotification(logId job.LogId, completedAt time.Time, output []byte, cmd_success bool) {
	cl, err := newFromConfig(configuration)
	if err != nil {
		logger.Error("can't create RPC client: %s", err.Error())
		return
	}
	defer cl.close()

	arg := &job.CompleteNotification{
		LogId:       logId,
		CompletedAt: completedAt,
		Output:      output,
		Success:     cmd_success,
	}
	var success bool
	cl.conn.SetDeadline(time.Now().Add(CONN_TIMEOUT))
	err = cl.Call("Server.CompleteNotification", arg, &success)
	if err != nil {
		logger.Error("can't send CompleteNotification message: %s", err.Error())
		return
	}

	if !success {
		logger.Error("server can't handle CompleteNotification: success is false")
	}
}
