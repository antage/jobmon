package logger

import (
	"fmt"
	"log/syslog"
	"os"
)

var logger *syslog.Writer

func Init(prefix string) {
	var err error
	logger, err = syslog.New(syslog.LOG_INFO, prefix)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't connect to syslog")
		os.Exit(3)
	}
}

func Close() {
	err := logger.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't close syslog")
	}
}

func Error(format string, a ...interface{}) {
	err := logger.Err(fmt.Sprintf(format, a...))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't send message to syslog")
		os.Exit(3)
	}
}

func Info(format string, a ...interface{}) {
	err := logger.Info(fmt.Sprintf(format, a...))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't send message to syslog")
		os.Exit(3)
	}
}
