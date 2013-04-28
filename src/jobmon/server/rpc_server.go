package main

import (
	"fmt"
	"jobmon/logger"
	"net"
	"net/rpc"
)

func runRpcServer(serverInstace *Server) error {
	rpcsrv := rpc.NewServer()
	err := rpcsrv.Register(serverInstace)
	if err != nil {
		logger.Error("can't register RPC methods: %s", err)
		return err
	}

	laddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", config.RPC.Listen, config.RPC.Port))
	if err != nil {
		logger.Error("can't resolve listen address for RPC: %s", err.Error())
		return err
	}
	l, err := net.ListenTCP("tcp4", laddr)
	if err != nil {
		logger.Error("can't bind to IPv4 tcp socket: %s", err)
		return err
	}

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				logger.Error("can't accept new connection for RPC: %s", err.Error())
			}

			if c != nil {
				if ipaddr, ok := c.RemoteAddr().(*net.TCPAddr); ok {
					allowed := false
					for _, allow := range config.RPC.allowIPNet {
						if allow.Contains(ipaddr.IP) {
							allowed = true
							break
						}
					}

					if allowed {
						go rpcsrv.ServeConn(c)
					} else {
						logger.Error("can't accept connection because remote address isn't allowed: %s", ipaddr.String())
						err := c.Close()
						if err != nil {
							logger.Error("can't close connection to RPC: %s", err.Error())
						}
					}
				} else {
					logger.Error("can't get remote address of accepted connection")
					err := c.Close()
					if err != nil {
						logger.Error("can't close connection to RPC: %s", err.Error())
					}
				}
			}
		}
	}()

	return nil
}
