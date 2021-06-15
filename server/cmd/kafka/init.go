package main

import (
	"os"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	log "github.com/sirupsen/logrus"
)

func initLog() {
	log.SetOutput(os.Stdout)
	log.SetLevel(config.LogLevel())
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	/*
		if graylogServer := config.GraylogServer(); graylogServer != "" {
			hook := graylog.NewGraylogHook(graylogServer, map[string]interface{}{"app": app, "version": version, "facility": app})
			defer hook.Flush()

			log.AddHook(hook)
		}
	*/
}
