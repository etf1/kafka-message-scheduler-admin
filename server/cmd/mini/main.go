package main

import (
	"net/http"
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
	version        = "mini"
	enable_metrics = false
)

func main() {
	if enable_metrics {
		metrics.DefaultConfig.CollectionInterval = time.Second
		if err := metrics.RunCollector(metrics.DefaultConfig); err != nil {
			log.Errorf("metrics error: %v", err)
		}
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
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
			kafkaRunner.Close()
		case <-exitchan:
			log.Printf("scheduler exited")
			break loop
		}
	}

	log.Printf("exiting ...")
}
