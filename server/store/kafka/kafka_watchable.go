package kafka

import (
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler/schedule/kafka"
	log "github.com/sirupsen/logrus"
)

type WatchableStore struct {
	consumers map[string]consumer
	processor
}

func NewWatchableStore(buckets []Bucket) (WatchableStore, error) {
	p := newProcessor(nil)
	p.start()

	ws := WatchableStore{
		consumers: make(map[string]consumer),
		processor: p,
	}

	for _, bucket := range buckets {
		c, err := newConsumer(bucket.Name, bucket.BootstrapServers, bucket.Topics)
		if err != nil {
			log.Errorf("cannot create kafka consumer for %+v: %v", bucket, err)
			continue
		}

		err = c.start(p.processChan)
		if err != nil {
			log.Errorf("cannot start kafka consumer for %+v: %v", bucket, err)
			continue
		}

		ws.consumers[bucket.Name] = c
	}

	return ws, nil
}

func (ws WatchableStore) Close() {
	defer log.Printf("kafka store closed")

	ws.processor.close()

	log.Printf("closing kafka store ...")
	for _, c := range ws.consumers {
		c.close()
	}
}

func (ws WatchableStore) Watch() (chan store.Event, error) {
	resultChan := make(chan store.Event, 10000)

	go func() {
		for e := range ws.processedChan {
			eventType := store.UpsertType
			if len(e.Value) == 0 {
				eventType = store.DeletedType
			}

			resultChan <- store.Event{
				EventType: eventType,
				Schedule: store.Schedule{
					SchedulerName: e.name,
					Schedule: kafka.Schedule{
						Message: e.Message,
					},
				},
			}
		}
	}()

	return resultChan, nil
}
