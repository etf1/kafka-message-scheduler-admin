package kafka

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/db/blevedb"
	"github.com/etf1/kafka-message-scheduler-admin/server/db/simple"
	"github.com/etf1/kafka-message-scheduler-admin/server/decoder"
	"github.com/etf1/kafka-message-scheduler-admin/server/decoder/httpdecoder"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/bbolt"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/rest"
)

const (
	FileMode = 0755
)

type Runner struct {
	stopChan chan bool
	exitChan chan bool
	dataDir  string
}

func (r Runner) Close() {
	r.stopChan <- true
	<-r.exitChan
}

func NewRunner(dataDir string) *Runner {
	return &Runner{
		stopChan: make(chan bool),
		exitChan: make(chan bool, 1),
		dataDir:  dataDir,
	}
}

func (r *Runner) Start() error {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovering from panic in runner: %v", r)
		}
		r.exitChan <- true
		log.Printf("kafka runner stopped")
	}()

	dir := r.dataDir
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	err := os.MkdirAll(dir, FileMode)
	if err != nil {
		return fmt.Errorf("cannot directories %v: %w", dir, err)
	}

	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		log.Println(variable[0], "=>", variable[1])
	}

	var dec decoder.Decoder

	if decoderURL := config.KafkaMessageBodyDecoder(); decoderURL != "" {
		dec = httpdecoder.Decoder{
			URL: decoderURL,
		}
	}

	// cold DB
	resolver := httpresolver.NewResolver(config.SchedulersAddr())
	bboltStore, err := bbolt.NewStore(dir + "schedules.bbolt")
	if err != nil {
		return fmt.Errorf("cannot create bbolt store: %w", err)
	}
	defer bboltStore.Close()

	watchableStore, err := NewWatchableStoreFromResolver(resolver, SchedulesTopics, dec)
	if err != nil {
		return fmt.Errorf("cannot create watchable store: %w", err)
	}
	defer watchableStore.Close()

	coldDB, err := blevedb.NewDB(blevedb.Config{
		InternalStore: bboltStore,
		SourceStore:   watchableStore,
		Path:          dir + "schedules.bleve",
	})
	if err != nil {
		return fmt.Errorf("cannot create bleve db %v: %w", dir, err)
	}
	defer coldDB.Close()

	// history DB
	historyBboltStore, err := bbolt.NewStore(dir + "history.bbolt")
	if err != nil {
		return fmt.Errorf("cannot create history bbolt store: %w", err)
	}
	defer historyBboltStore.Close()

	historyWatchableStore, err := NewWatchableStoreFromResolver(resolver, HistoryTopic, dec)
	if err != nil {
		return fmt.Errorf("cannot create history watchable store: %w", err)
	}
	defer historyWatchableStore.Close()

	historyDB, err := blevedb.NewDB(blevedb.Config{
		InternalStore: historyBboltStore,
		SourceStore:   historyWatchableStore,
		Path:          dir + "history.bleve",
	})
	if err != nil {
		return fmt.Errorf("cannot create history bleve db: %w", err)
	}
	defer historyDB.Close()

	// live DB
	liveDB := simple.DB{
		Store: rest.NewStore(resolver, dec),
	}

	srv := runner.NewServer(coldDB, liveDB, historyDB, resolver)

	helper.StartupHTTPServer(srv)
	<-r.stopChan
	helper.LogErr(helper.ShutdownHTTPServer(srv))

	return nil
}
