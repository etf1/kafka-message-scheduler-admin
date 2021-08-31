package kafka

import (
	"fmt"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/kafka"
	log "github.com/sirupsen/logrus"
)

var (
	DefaultTopics = SchedulesTopics
)

func SchedulesTopics(s httpresolver.Scheduler) []string {
	return s.Topics()
}

func HistoryTopic(s httpresolver.Scheduler) []string {
	return []string{s.History()}
}

type WatchableStoreFromResolver struct {
	*kafka.WatchableStore
	resolver httpresolver.Resolver
	schs     []schedulers.Scheduler
	topics   func(s httpresolver.Scheduler) []string
	// used by close
	stopChan chan bool
	exitChan chan bool
}

func (wr *WatchableStoreFromResolver) Close() {
	log.Printf("calling WatchableStoreFromResolver close")
	wr.stopChan <- true
	<-wr.exitChan
}

func (wr *WatchableStoreFromResolver) updateBuckets() error {
	schs, err := wr.resolver.List()
	if err != nil {
		return err
	}

	log.Warnf("retrieved schedulers=%+v", schs)

	wr.schs = schs
	buckets := make([]kafka.Bucket, 0)

	for _, sch := range schs {
		s, ok := sch.(httpresolver.Scheduler)
		if !ok {
			return fmt.Errorf("unable to cast: %T", sch)
		}
		buckets = append(buckets, kafka.Bucket{
			Name:             s.Name(),
			BootstrapServers: s.BootstrapServers(),
			Topics:           wr.topics(s),
		})
	}

	wr.WatchableStore.AddBuckets(buckets...)

	return nil
}

type TopicFunc func(s httpresolver.Scheduler) []string

func NewWatchableStoreFromResolver(resolver httpresolver.Resolver, topics TopicFunc) (*WatchableStoreFromResolver, error) {
	wr := &WatchableStoreFromResolver{
		resolver: resolver,
		stopChan: make(chan bool, 1),
		exitChan: make(chan bool, 1),
		topics:   topics,
	}

	ws, err := kafka.NewWatchableStore()
	if err != nil {
		return wr, err
	}

	wr.WatchableStore = ws

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("recovering from panic in watchable store resolver: %v", r)
			}
			wr.exitChan <- true
			log.Warnf("watchable store resolver exited.")
		}()

		duration := 30 * time.Second
		recheck := time.NewTimer(duration)
		defer recheck.Stop()

		// initial update of the buckets
		err := wr.updateBuckets()
		if err != nil {
			log.Errorf("unable to update buckets: %v", err)
		}

	loop:
		for {
			// to avoid <-time.After(..) memory leak
			recheck.Reset(duration)

			select {
			case <-recheck.C:
				// after a period, update buckets
				err := wr.updateBuckets()
				if err != nil {
					log.Errorf("unable to update buckets: %v", err)
				}
			case <-wr.stopChan:
				wr.WatchableStore.Close()
				break loop
			}
		}
	}()

	return wr, nil
}

/*
func NewWatchableStoreFromResolver(r httpresolver.Resolver, topics func(s httpresolver.Scheduler) []string) (*kafka.WatchableStore, error) {
	lst, err := r.List()
	if err != nil {
		return nil, err
	}

	buckets := make([]kafka.Bucket, 0)

	for _, sch := range lst {
		s, ok := sch.(httpresolver.Scheduler)
		if !ok {
			log.Errorf("unable to cast: %T", sch)
			continue
		}
		buckets = append(buckets, kafka.Bucket{
			Name:             s.Name(),
			BootstrapServers: s.BootstrapServers(),
			Topics:           topics(s),
		})
	}

	return kafka.NewWatchableStore(buckets...)
}
*/
