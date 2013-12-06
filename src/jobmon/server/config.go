package main

import (
	"encoding/json"
	"fmt"
	"io"
	"jobmon/logger"
	"net"
	"os"
)

const CONF_DIR = "/etc/jobmon"

const ASSETS_DIR = "/usr/share/jobmond/assets"

type webConfig struct {
	Listen string
	Port   uint
}

type rpcConfig struct {
	Listen     string
	Port       uint
	Allow      []string
	allowIPNet []*net.IPNet
}

type smtpConfig struct {
	Host string
	Port uint
}

type messageConfig struct {
	From    string
	Subject string
}

type mailConfig struct {
	SMTP    smtpConfig
	Message messageConfig
	Notify  []string
}

type mainConfig struct {
	Web  webConfig
	RPC  rpcConfig
	Mail mailConfig
}

var config *mainConfig

func defaultConfig() *mainConfig {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Error("can't get hostname: %s", err.Error())
		hostname = "localhost"
	}

	c := &mainConfig{
		Web: webConfig{
			Listen: "127.0.0.1",
			Port:   8080,
		},
		RPC: rpcConfig{
			Listen:     "127.0.0.1",
			Port:       10432,
			Allow:      []string{"127.0.0.0/8"},
			allowIPNet: nil,
		},
		Mail: mailConfig{
			SMTP: smtpConfig{
				Host: "127.0.0.1",
				Port: 25,
			},
			Message: messageConfig{
				From:    fmt.Sprintf("no-reply@%s", hostname),
				Subject: "[JOBMOND] Job fail notification",
			},
			Notify: []string{"root"},
		},
	}

	c.RPC.allowIPNet = make([]*net.IPNet, 0, len(c.RPC.Allow))
	for _, allow := range c.RPC.Allow {
		_, ipnet, err := net.ParseCIDR(allow)
		if err != nil {
			logger.Error("can't parse allow CIDR (%s) in default configuration: %s", err.Error())
			continue
		}
		c.RPC.allowIPNet = append(c.RPC.allowIPNet, ipnet)
	}

	return c
}

func parseConfigFile(filename string) error {
	configFile, err := os.Open(filename)
	if err != nil {
		logger.Error("can't open config file: %s", err.Error())
		return err
	}

	decoder := json.NewDecoder(configFile)
	config = defaultConfig()
	err = decoder.Decode(config)
	if err != nil && err != io.EOF {
		logger.Error("can't parse config file: %s", err.Error())
		return err
	}

	config.RPC.allowIPNet = make([]*net.IPNet, 0, len(config.RPC.Allow))
	for _, allow := range config.RPC.Allow {
		_, ipnet, err := net.ParseCIDR(allow)
		if err != nil {
			logger.Error("can't parse allow CIDR (%s) in loaded configuration: %s", err.Error())
			continue
		}
		config.RPC.allowIPNet = append(config.RPC.allowIPNet, ipnet)
	}

	return nil
}
