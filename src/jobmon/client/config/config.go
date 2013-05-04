package config

import (
	"encoding/json"
	"io"
	"jobmon/logger"
	"os"
)

const CONF_DIR = "/etc/jobmon"

type RPC struct {
	Host string
	Port uint
}

type Main struct {
	RPC RPC
}

func Default() *Main {
	return &Main{
		RPC: RPC{
			Host: "127.0.0.1",
			Port: 10432,
		},
	}
}

func ParseFile(filename string) (*Main, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		logger.Error("can't open config file: %s", err.Error())
		return nil, err
	}

	decoder := json.NewDecoder(configFile)
	config := Default()
	err = decoder.Decode(config)
	if err != nil && err != io.EOF {
		logger.Error("can't parse config file: %s", err.Error())
		return nil, err
	}

	return config, nil
}
