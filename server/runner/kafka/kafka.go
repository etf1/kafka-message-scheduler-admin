package kafka

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/db/blevedb"
	"github.com/etf1/kafka-message-scheduler-admin/server/db/simple"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
	"github.com/etf1/kafka-message-scheduler-admin/server/restapi"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/bbolt"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/kafka"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/rest"
)

type Runner struct {
	stopChan chan bool
	dataDir  string
}

func (r Runner) Close() {
	r.stopChan <- true
}

func NewRunner(dataDir string) *Runner {
	return &Runner{
		stopChan: make(chan bool, 1),
		dataDir:  dataDir,
	}
}

func (r *Runner) Start() error {
	dir := r.dataDir
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("cannot directories %v: %w", dir, err)
	}

	resolver := httpresolver.NewResolver(config.SchedulersAddr())
	bboltStore, err := bbolt.NewStore(dir + "schedules.bbolt")
	if err != nil {
		return fmt.Errorf("cannot create bbolt store: %w", err)
	}

	watchableStore, err := kafka.NewWatchableStoreFromResolver(resolver)
	if err != nil {
		return fmt.Errorf("cannot create watchable store: %w", err)
	}

	cold, err := blevedb.NewDB(blevedb.Config{
		InternalStore: bboltStore,
		SourceStore:   watchableStore,
		Path:          dir + "schedules.bleve",
	})
	if err != nil {
		return fmt.Errorf("cannot create bleve db: %w", err)
	}

	live := simple.DB{
		Store: rest.NewStore(resolver),
	}

	//srv := restapi.NewServer(coldb, livedb, resolver)
	//srv.Start(config.ServerAddr())

	if config.APIServerOnly() {
		srv := restapi.NewServer(simple.DB{
			Store: cold,
		}, simple.DB{
			Store: live,
		}, resolver)
		srv.Start(config.ServerAddr())
		defer srv.Stop()
	} else {
		mainRouter := mux.NewRouter().StrictSlash(true)

		mainRouter.PathPrefix("/api").Handler(http.StripPrefix("/api", restapi.NewRouter(simple.DB{
			Store: cold,
		}, simple.DB{
			Store: live,
		}, resolver)))

		mainRouter.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(config.StaticFilesDir()))))

		addr := config.ServerAddr()
		srv := &http.Server{
			Handler:      mainRouter,
			Addr:         addr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		defer helper.ShutdownHttpServer(srv, config.ShutdownTimeout())

		go func() {
			log.Printf("starting server on %s", addr)
			defer log.Printf("server stopped")
			if err := srv.ListenAndServe(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	<-r.stopChan

	return nil
}
