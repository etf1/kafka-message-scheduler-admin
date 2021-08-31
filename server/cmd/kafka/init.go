package main

import (
	"context"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	graylog "github.com/gemnasium/logrus-graylog-hook"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func initLog() {
	log.SetOutput(os.Stdout)
	log.SetLevel(config.LogLevel())
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	if graylogServer := config.GraylogServer(); graylogServer != "" {
		hook := graylog.NewGraylogHook(graylogServer, map[string]interface{}{"app": app, "version": version, "facility": app})
		defer hook.Flush()

		log.AddHook(hook)
	}
}

func initPprof(enabledBydefault bool) func() {
	timeout := 5 * time.Second
	startchan := make(chan os.Signal, 1)
	signal.Notify(startchan, syscall.SIGUSR1)

	stopchan := make(chan os.Signal, 1)
	signal.Notify(stopchan, syscall.SIGUSR2)

	exitchan := make(chan bool)

	var server *http.Server

	shutdown := func() {
		if server != nil {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			log.Printf("shutting down pprof server")
			err := server.Shutdown(ctx)
			if err != nil {
				log.Printf("failed to stop pprof server: %v", err)
			}
		}
	}

	closePprof := func() {
		shutdown()
		exitchan <- true
	}

	router := http.NewServeMux()
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	go func() {
		log.Warnf("pprof launcher started")
		defer log.Printf("pprof launcher exited")

		for {
			select {
			case <-exitchan:
				return
			case <-stopchan:
				shutdown()
			case <-startchan:
				server = &http.Server{
					Addr:    "localhost:6060",
					Handler: router,
				}
				go func() {
					log.Warnf("starting http pprof server")
					log.Println(server.ListenAndServe())
					log.Warnf("http server pprof shutted down")
					server = nil
				}()
			}
		}
	}()
	if enabledBydefault {
		startchan <- syscall.SIGUSR1
	}
	return closePprof
}

func initProm() func() {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	srv := &http.Server{
		Addr:    config.MetricsAddr(),
		Handler: mux,
	}

	exitChan := make(chan bool)
	go func() {
		log.Printf("prometheus metrics available on %s/metrics", srv.Addr)
		defer log.Printf("prometheus metrics server stopped")

		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Error(err)
			}
		}
		exitChan <- true
	}()

	return func() {
		err := helper.ShutdownHTTPServer(srv)
		if err != nil {
			log.Errorf("cannot close prometheus server: %v", err)
			return
		}
		<-exitChan
	}
}
