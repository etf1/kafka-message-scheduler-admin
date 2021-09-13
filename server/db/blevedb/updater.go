package blevedb

import (
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler/schedule"
	log "github.com/sirupsen/logrus"
)

const (
	MaxChanSize      = 10000
	MaxBatchChanSize = 1000
)

type updater struct {
	input chan event
	store.BatchableStore
}

func newUpdater(bs store.BatchableStore) updater {
	return updater{
		make(chan event, MaxChanSize),
		bs,
	}
}

func (u updater) start() {
	defer log.Printf("updater closed")

	batchChan := make(chan store.Event, MaxBatchChanSize)
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
				log.Debugf("batch index: %T %+v ", sch, sch)
				batchChan <- store.Event{
					EventType: store.UpsertType,
					Schedule:  sch,
				}
			case deleteType:
				log.Debugf("batch delete: %T %v", sch, sch)
				batchChan <- store.Event{
					EventType: store.DeletedType,
					Schedule:  sch,
				}
			}
		}
	}
}

func (u updater) upsert(id string, s schedule.Schedule) {
	if u.input == nil {
		return
	}
	u.input <- event{
		eventType: upsertType,
		id:        id,
		data:      s,
	}
}

func (u updater) delete(id string, s schedule.Schedule) {
	if u.input == nil {
		return
	}
	u.input <- event{
		eventType: deleteType,
		id:        id,
		data:      s,
	}
}
