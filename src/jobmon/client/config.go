package main

import (
	"encoding/json"
	"io"
	"jobmon/logger"
	"os"
)

const CONF_DIR = "/etc/jobmon"

type rpcConfig struct {
	Host string
	Port uint
}

type mainConfig struct {
	RPC rpcConfig
}

var config *mainConfig

func defaultConfig() *mainConfig {
	c := &mainConfig{
		RPC: rpcConfig{
			Host: "127.0.0.1",
			Port: 10432,
		},
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

	return nil
}
