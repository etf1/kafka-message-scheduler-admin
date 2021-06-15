package blevedb

// import (
// 	"fmt"
// 	"time"

// 	"github.com/etf1/kafka-message-scheduler-admin/server/store"
// 	"github.com/etf1/kafka-message-scheduler/schedule"
// 	log "github.com/sirupsen/logrus"
// )

// type updater struct {
// 	input chan event
// 	store.MutableStore
// }

// func newUpdater(s store.MutableStore) (updater, error) {
// 	return updater{
// 		make(chan event, 1000000),
// 		s,
// 	}, nil
// }

// func (u updater) start() error {
// 	defer log.Printf("indexer closed")

// 	duration := 500 * time.Millisecond
// 	timeout := time.NewTimer(duration)
// 	defer timeout.Stop()

// 	counter := 0

// 	batch := make([]event, 0)

// 	processBatch := func() {
// 		log.Printf("batch indexing %v documents", counter)

// 		batch = make([]event, 0)
// 	}

// loop:
// 	for {
// 		timeout.Reset(duration)
// 		select {
// 		case evt, ok := <-u.input:
// 			log.Printf("received event from input channel")
// 			if !ok {
// 				log.Printf("input channel closed")
// 				processBatch()
// 				break loop
// 			}

// 			if len(batch)%batch_size == 0 {
// 				processBatch()
// 			}

// 			/*
// 				switch evt.eventType {
// 				case upsertType:
// 					log.Printf("batch index: %+v", evt.doc)
// 					err := batch.Index(evt.id, evt.doc)
// 					if err != nil {
// 						log.Errorf("index batch failed: %v", err)
// 						break
// 					}
// 				case deleteType:
// 					log.Printf("batch delete: %v", evt.doc)
// 					batch.Delete(evt.id)
// 				}
// 				counter++
// 				if counter%batch_size == 0 {
// 					indexBatch()
// 					//batch = i.NewBatch()
// 					log.Warnf("indexed %v documents", counter)
// 				}
// 			*/
// 		case <-timeout.C:
// 			log.Warnf("input channel timeout")
// 			if batch.Size() != 0 {
// 				indexBatch()
// 				log.Warnf("indexed %v documents", counter)
// 			}
// 		}
// 	}

// }

// func (u updater) upsertData(id string, s schedule.Schedule) error {
// 	if u.input == nil {
// 		return fmt.Errorf("indexer not initialized or closed")
// 	}
// 	u.input <- event{
// 		upsertType,
// 		id,
// 		s,
// 	}
// 	return nil
// }

// func (u updater) deleteData(id string, s schedule.Schedule) error {
// 	if u.input == nil {
// 		return fmt.Errorf("indexer not initialized or closed")
// 	}
// 	u.input <- event{
// 		deleteType,
// 		id,
// 		s,
// 	}
// 	return nil
// }
