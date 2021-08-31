package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/etf1/kafka-message-scheduler-admin/server/runner/mini"
	log "github.com/sirupsen/logrus"
	metrics "github.com/tevjef/go-runtime-metrics"
)

var (
	version             = "mini"
	enableTevjefMetrics = false
)

func main() {
	if enableTevjefMetrics {
		metrics.DefaultConfig.CollectionInterval = time.Second
		if err := metrics.RunCollector(metrics.DefaultConfig); err != nil {
			log.Errorf("metrics error: %v", err)
		}
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	kafkaRunner := mini.NewRunner()

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
