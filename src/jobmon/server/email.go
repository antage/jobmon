package main

import (
	"bytes"
	"fmt"
	"jobmon/job"
	"jobmon/logger"
	"net/smtp"
)

func sendNotification(server *Server, logEntry *job.LogEntry) {
	logger.Info("send notification: log entry id=%d", logEntry.Id)

	c, err := smtp.Dial("localhost:25")
	if err != nil {
		logger.Error("can't dial smtp server: %s", err.Error())
		return
	}

	for _, recipient := range config.Mail.Notify {

		err = c.Mail(config.Mail.Message.From)
		if err != nil {
			logger.Error("can't set sender for smtp: %s", err.Error())
			return
		}

		err = c.Rcpt(recipient)
		if err != nil {
			logger.Error("can't set recipient for smtp: %s", err.Error())
			return
		}

		wc, err := c.Data()
		if err != nil {
			logger.Error("can't get message data: %s", err.Error())
			return
		}
		defer wc.Close()

		data := bytes.NewBuffer([]byte{})
		data.WriteString(fmt.Sprintf("From: %s\n", config.Mail.Message.From))
		data.WriteString(fmt.Sprintf("To: %s\n", recipient))
		data.WriteString(fmt.Sprintf("Subject: %s\n", config.Mail.Message.Subject))
		data.WriteString("\n")
		data.WriteString(fmt.Sprintf("Log id=%d\n", logEntry.Id))
		data.Write(logEntry.Output)
		// fill data
		if _, err = data.WriteTo(wc); err != nil {
			logger.Error("can't write data in smtp: %s", err.Error())
			return
		}
	}
}
