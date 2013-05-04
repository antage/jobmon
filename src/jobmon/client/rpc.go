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
const IDLE_TIMEOUT = 600 * time.Second

type rpcClient struct {
	*rpc.Client
	configuration *config.Main
	conn          net.Conn
	closed        bool
}

func newFromConfig(c *config.Main) (*rpcClient, error) {
	conn, err := net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", c.RPC.Host, c.RPC.Port), CONN_TIMEOUT)
	if err != nil {
		logger.Error("can't resolve connect to server: %s", err.Error())
		return nil, err
	}
	if tcpconn, ok := conn.(*net.TCPConn); ok {
		tcpconn.SetKeepAlive(true)
	}

	return &rpcClient{rpc.NewClient(conn), c, conn, false}, nil
}

func (cl *rpcClient) close() error {
	if !cl.closed {
		err := cl.Close()
		if err != nil {
			logger.Error("can't close RPC client: %s", err.Error())
			return err
		}
		err = cl.conn.Close()
		if err != nil {
			logger.Error("can't close connection: %s", err.Error())
			return err
		}
		cl.closed = true
	}
	return nil
}

func (cl *rpcClient) startNotification(jobId job.JobId) (logId job.LogId, err error) {
	arg := &job.StartNotification{
		JobId:     &jobId,
		StartedAt: time.Now(),
	}

	cl.conn.SetDeadline(time.Now().Add(CONN_TIMEOUT))
	err = cl.Call("Server.StartNotification", arg, &logId)
	cl.conn.SetDeadline(time.Now().Add(IDLE_TIMEOUT))

	if err != nil {
		logger.Error("can't send StartNotification message: %s", err.Error())
		return job.LogId(0), err
	}

	return logId, nil
}

func (cl *rpcClient) aliveNotification(logId job.LogId) {
	arg := &job.AliveNotification{LogId: logId}

	var success bool

	cl.conn.SetDeadline(time.Now().Add(CONN_TIMEOUT))
	err := cl.Call("Server.AliveNotification", arg, &success)
	cl.conn.SetDeadline(time.Now().Add(IDLE_TIMEOUT))

	if err != nil {
		logger.Error("can't send AliveNotification message: %s", err.Error())
		return
	}

	if !success {
		logger.Error("server can't handle AliveNotification: success is false")
	}
}

func (cl *rpcClient) completeNotification(logId job.LogId, completedAt time.Time, output []byte, cmd_success bool) {
	arg := &job.CompleteNotification{
		LogId:       logId,
		CompletedAt: completedAt,
		Output:      output,
		Success:     cmd_success,
	}
	var success bool
	cl.conn.SetDeadline(time.Now().Add(CONN_TIMEOUT))
	err := cl.Call("Server.CompleteNotification", arg, &success)
	cl.conn.SetDeadline(time.Now().Add(IDLE_TIMEOUT))
	if err != nil {
		logger.Error("can't send CompleteNotification message: %s", err.Error())
		return
	}
	if !success {
		logger.Error("server can't handle CompleteNotification: success is false")
	}
}
