package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner/kafka"
	log "github.com/sirupsen/logrus"
	metrics "github.com/tevjef/go-runtime-metrics"
)

var (
	app                 = "kafka-message-scheduler-admin"
	version             = "undefined"
	enableTevjefMetrics = true
)

func main() {
	initLog()

	if enableTevjefMetrics {
		metrics.DefaultConfig.CollectionInterval = time.Second
		if err := metrics.RunCollector(metrics.DefaultConfig); err != nil {
			log.Errorf("metrics error: %v", err)
		}
	}

	defer initProm()()
	defer initPprof(true)()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	kafkaRunner := kafka.NewRunner(config.DataRootDir())

	exitchan := make(chan bool)

	go func() {
		log.Printf("starting version=%v", version)
		if err := kafkaRunner.Start(); err != nil {
			log.Errorf("failed to start runner: %v", err)
		}
		exitchan <- true
	}()

loop:
	for {
		select {
		case <-sigchan:
			log.Printf("closing runner...")
			kafkaRunner.Close()
		case <-exitchan:
			log.Printf("runner exited")
			break loop
		}
	}

	log.Printf("done.")
}
