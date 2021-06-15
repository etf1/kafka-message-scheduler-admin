package kafka

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/db/blevedb"
	"github.com/etf1/kafka-message-scheduler-admin/server/db/simple"
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

	coldb, err := blevedb.NewDB(blevedb.Config{
		InternalStore: bboltStore,
		SourceStore:   watchableStore,
		Path:          dir + "schedules.bleve",
	})
	if err != nil {
		return fmt.Errorf("cannot create bleve db: %w", err)
	}

	// coldb := simple.DB{
	// 	Store: hmap.NewStore(),
	// }

	livedb := simple.DB{
		Store: rest.NewStore(resolver),
	}

	srv := restapi.NewServer(coldb, livedb, resolver)
	srv.Start(config.APIServerAddr())

	<-r.stopChan

	srv.Stop()
	log.Printf("kafka runner closed")

	return nil
}
