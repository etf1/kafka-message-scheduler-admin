package kafka

import (
	"fmt"
	"reflect"

	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler/schedule/kafka"
	log "github.com/sirupsen/logrus"
)

type WatchableStore struct {
	consumers map[string]consumer
	processor
}

func NewWatchableStore(buckets ...Bucket) (*WatchableStore, error) {
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

	return &ws, nil
}

func (ws WatchableStore) Close() {
	defer log.Warnf("watchable kafka store closed")

	ws.processor.close()

	log.Warnf("closing watchable kafka store ...")
	for _, c := range ws.consumers {
		c.close()
	}
}

func (ws *WatchableStore) AddBuckets(buckets ...Bucket) {
	for _, bucket := range buckets {
		c, found := ws.consumers[bucket.Name]

		// consumer config changed
		if found && (c.bootstrapServers != bucket.BootstrapServers || !reflect.DeepEqual(c.topics, bucket.Topics)) {
			// so sending reset event
			ws.processedChan <- event{
				storeResetType,
				bucket.Name,
				nil,
			}
			// closing current consumer
			c.close()
		} else if found {
			// nothing changed
			continue
		}
		fmt.Printf("setting new consumer %v\n", bucket)
		// starting new consumer
		c, err := newConsumer(bucket.Name, bucket.BootstrapServers, bucket.Topics)
		if err != nil {
			log.Errorf("cannot create kafka consumer for %+v: %v", bucket, err)
			continue
		}

		err = c.start(ws.processChan)
		if err != nil {
			log.Errorf("cannot start kafka consumer for %+v: %v", bucket, err)
			continue
		}

		ws.consumers[bucket.Name] = c
	}
}

func (ws WatchableStore) Watch() (chan store.Event, error) {
	resultChan := make(chan store.Event, ChanSize)

	go func() {
		for e := range ws.processedChan {
			switch e.evtType {
			case messageType:
				evtType := store.UpsertType
				if len(e.Value) == 0 {
					evtType = store.DeletedType
				}
				resultChan <- store.Event{
					EventType: evtType,
					Schedule: store.Schedule{
						SchedulerName: e.name,
						Schedule: kafka.Schedule{
							Message: e.Message,
						},
					},
				}
			case storeResetType:
				resultChan <- store.Event{
					EventType: store.StoreResetType,
				}
			}
		}
	}()

	return resultChan, nil
}
