package blevedb

import (
	"fmt"

	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler/schedule"
	log "github.com/sirupsen/logrus"
)

type updater struct {
	input chan event
	store.BatchableStore
}

func newUpdater(bs store.BatchableStore) updater {
	return updater{
		make(chan event, 10000),
		bs,
	}
}

func (u updater) start() error {
	defer log.Printf("updater closed")

	batchChan := make(chan store.Event, 1000)
	defer close(batchChan)

	errChan := u.Batch(batchChan)

loop:
	for {
		select {
		case err, ok := <-errChan:
			if !ok {
				break loop
			}
			log.Errorf("received error from batch: %v", err)
		case evt, ok := <-u.input:
			if !ok {
				log.Printf("input channel closed")
				break loop
			}
			sch, ok := evt.data.(store.Schedule)
			if !ok {
				log.Errorf("unexpected schedule object: %T", evt.data)
			}
			switch evt.eventType {
			case upsertType:
				log.Debugf("batch index: %+v", sch)
				batchChan <- store.Event{
					EventType: store.UpsertType,
					Schedule:  sch,
				}
			case deleteType:
				log.Debugf("batch delete: %v", sch)
				batchChan <- store.Event{
					EventType: store.DeletedType,
					Schedule:  sch,
				}
			}
		}
	}
	return nil
}

func (u updater) upsertData(id string, s schedule.Schedule) error {
	if u.input == nil {
		return fmt.Errorf("indexer not initialized or closed")
	}
	u.input <- event{
		upsertType,
		id,
		s,
	}
	return nil
}

func (u updater) deleteData(id string, s schedule.Schedule) error {
	if u.input == nil {
		return fmt.Errorf("indexer not initialized or closed")
	}
	u.input <- event{
		deleteType,
		id,
		s,
	}
	return nil
}
