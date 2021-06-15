package blevedb

// import (
// 	"fmt"

// 	"github.com/etf1/kafka-message-scheduler-admin/server/store"
// 	"github.com/etf1/kafka-message-scheduler/schedule"
// 	log "github.com/sirupsen/logrus"
// )

// type updater struct {
// 	input chan event
// 	store.MutableStore
// }

// func newUpdater(s store.MutableStore) updater {
// 	return updater{
// 		make(chan event, 1000000),
// 		s,
// 	}
// }

// func (u updater) start() error {
// 	defer log.Printf("updater closed")
// 	/*
// 		duration := 500 * time.Millisecond
// 		timeout := time.NewTimer(duration)
// 		defer timeout.Stop()

// 		var batch []schedule.Schedule
// 	*/
// 	counter := 0
// loop:
// 	for {
// 		//timeout.Reset(duration)
// 		select {
// 		case evt, ok := <-u.input:
// 			log.Printf("updater: received event from input channel")
// 			if !ok {
// 				log.Printf("input channel closed")
// 				break loop
// 			}

// 			sch, ok := evt.data.(store.Schedule)
// 			if !ok {
// 				log.Errorf("unexpected schedule object: %T", evt.data)
// 			}

// 			//batch = append(batch, sch.Schedule)

// 			switch evt.eventType {
// 			case upsertType:
// 				log.Printf("batch index: %+v", sch)
// 				err := u.Add(sch.SchedulerName, sch.Schedule)
// 				if err != nil {
// 					log.Errorf("cannot add schedule %+v to store: %v", evt, err)
// 				}
// 			case deleteType:
// 				log.Printf("batch delete: %v", sch)
// 				err := u.Delete(sch.SchedulerName, sch.Schedule)
// 				if err != nil {
// 					log.Errorf("cannot add schedule %+v to store: %v", evt, err)
// 				}

// 			}

// 			counter++
// 			if counter%batch_size == 0 {
// 				log.Warnf("stored %v schedules", counter)
// 			}
// 			/*
// 				case <-timeout.C:
// 					log.Printf("input channel timeout")
// 					if batch.Size() != 0 {
// 						indexBatch()
// 						log.Warnf("indexed %v documents", counter)
// 					}
// 			*/
// 		}
// 	}
// 	return nil
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
