package bbolt

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler/schedule"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	batch_size = 1000
)

type DB struct {
	bbolt_db *bolt.DB
	// indexEpoch []indexEpochItem
}

type Schedule struct {
	ScheduleID        string `json:"id"`
	ScheduleEpoch     int64  `json:"epoch"`
	ScheduleTimestamp int64  `json:"timestamp"`
	Topic             string `json:"topic"`
	TargetTopic       string `json:"target-topic"`
	TargetKey         string `json:"target-key"`
	Value             []byte `json:"value"`
}

func NewSchedule(id, epoch interface{}, timestamp ...time.Time) Schedule {
	var sid string

	switch v := id.(type) {
	case int:
		sid = strconv.Itoa(v)
	case int64:
		sid = strconv.FormatInt(v, 10)
	case string:
		sid = v
	default:
		sid = ""
	}

	var iepoch int64
	switch v := epoch.(type) {
	case int:
		iepoch = int64(v)
	case int64:
		iepoch = v
	case time.Time:
		iepoch = v.Unix()
	default:
		iepoch = time.Now().Unix()
	}

	ts := time.Now().Unix()
	if len(timestamp) == 1 {
		ts = timestamp[0].Unix()
	}

	return Schedule{
		ScheduleID:        sid,
		ScheduleEpoch:     iepoch,
		ScheduleTimestamp: ts,
	}
}

func (s Schedule) ID() string {
	return s.ScheduleID
}

func (s Schedule) Epoch() int64 {
	return s.ScheduleEpoch
}

func (s Schedule) Timestamp() int64 {
	return s.ScheduleTimestamp
}

func (s Schedule) String() string {
	return fmt.Sprintf("{id:%s epoch:%v date:%v timestamp:%v}", s.ID(), s.Epoch(), time.Unix(s.Epoch(), 0), s.Timestamp())
}

func NewStore(path string) (DB, error) {
	db, err := bolt.Open(path, 0666, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return DB{}, err
	}
	db.MaxBatchSize = batch_size
	return DB{
		bbolt_db: db,
	}, nil
}

func (d DB) Get(schedulerName string, scheduleID string) ([]store.Schedule, error) {
	var schedules []store.Schedule

	err := d.bbolt_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(schedulerName))
		if b == nil {
			return nil
		}
		v := b.Get([]byte(scheduleID))
		if v != nil {
			var arr []Schedule
			err := json.Unmarshal(v, &arr)
			if err != nil {
				return err
			}

			schedules = make([]store.Schedule, len(arr))
			for i := 0; i < len(arr); i++ {
				schedules[i] = store.Schedule{
					SchedulerName: schedulerName,
					Schedule:      arr[i],
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (d DB) List(schedulerName string) (chan store.Schedule, error) {
	result := make(chan store.Schedule, 1000)

	go func() {
		defer close(result)

		d.bbolt_db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(schedulerName))
			if b == nil {
				return nil
			}
			b.ForEach(func(k, v []byte) error {
				if v != nil {
					var arr []Schedule
					err := json.Unmarshal(v, &arr)
					if err != nil {
						log.Errorf("unable to unmarshall: %v", err)
					} else {
						result <- store.Schedule{
							SchedulerName: schedulerName,
							Schedule:      arr[0],
						}
					}
				}
				return nil
			})
			return nil
		})
	}()

	return result, nil
}

func (d DB) Add(schedulerName string, ss ...schedule.Schedule) error {
	return d.bbolt_db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(schedulerName))
		if err != nil {
			return fmt.Errorf("cannot create bucket %s: %s", schedulerName, err)
		}
		for _, s := range ss {
			v := b.Get([]byte(s.ID()))
			if v == nil {
				buf, err := json.Marshal([]schedule.Schedule{s})
				if err != nil {
					log.Errorf("cannot marshall schedule %v: %v", s, err)
					continue
				}
				log.Debugf("bbolt store put (1) %s size %v", s.ID(), len(buf))
				err = b.Put([]byte(s.ID()), buf)
				if err != nil {
					log.Errorf("cannot put schedule %v: %v", s, err)
				}
				continue
			}

			buf, err := addToSlice(s, v)
			if err != nil {
				log.Errorf("cannot get bytes for %v: %v", string(v), err)
				continue
			}
			log.Debugf("bbolt store put (2) %s size %v", s.ID(), len(buf))
			err = b.Put([]byte(s.ID()), buf)
			if err != nil {
				log.Errorf("cannot put schedule %v: %v", s, err)
				continue
			}
		}

		return nil
	})
}

func (d DB) Close() {
	d.bbolt_db.Close()
}

func (d DB) Delete(schedulerName string, ss ...schedule.Schedule) error {
	return d.bbolt_db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(schedulerName))
		if b == nil {
			return nil
		}
		for _, s := range ss {
			err := b.Delete([]byte(s.ID()))
			if err != nil {
				log.Errorf("cannot delete schedule %v in bucket %v: %v", s, b, err)
			}
		}

		return nil
	})
}

func (d DB) removeFromBucket(tx *bolt.Tx, bucketName string, sch schedule.Schedule) error {
	b := tx.Bucket([]byte(bucketName))
	if b == nil {
		return nil
	}
	err := b.Delete([]byte(sch.ID()))
	if err != nil {
		log.Errorf("cannot delete schedule %v in bucket %v: %v", sch, b, err)
	}

	return nil
}

func (d DB) addToBucket(tx *bolt.Tx, bucketName string, sch schedule.Schedule) error {
	b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return fmt.Errorf("cannot create bucket %s: %s", bucketName, err)
	}
	v := b.Get([]byte(sch.ID()))
	if v == nil {
		buf, err := json.Marshal([]schedule.Schedule{sch})
		if err != nil {
			return fmt.Errorf("cannot marshall schedule %v: %v", sch, err)
		}
		log.Printf("bbolt store put (1) %s size %v", sch.ID(), len(buf))
		err = b.Put([]byte(sch.ID()), buf)
		if err != nil {
			return fmt.Errorf("cannot put schedule %v: %v", sch, err)
		}
		return nil
	}

	buf, err := addToSlice(sch, v)
	if err != nil {
		return fmt.Errorf("cannot get bytes for %v: %v", string(v), err)
	}
	log.Printf("bbolt store put (2) %s size %v", sch.ID(), len(buf))
	err = b.Put([]byte(sch.ID()), buf)
	if err != nil {
		return fmt.Errorf("cannot put schedule %v: %v", sch, err)
	}

	return nil
}

func (d DB) Batch(events chan store.Event) chan error {
	errChan := make(chan error, batch_size)

	var batch []store.Event

	processBatch := func() error {
		return d.bbolt_db.Batch(func(tx *bolt.Tx) error {
			for _, evt := range batch {
				var err error
				switch evt.EventType {
				case store.DeletedType:
					err = d.removeFromBucket(tx, evt.SchedulerName, evt.Schedule.Schedule)
				case store.UpsertType:
					err = d.addToBucket(tx, evt.SchedulerName, evt.Schedule.Schedule)
				}
				if err != nil {
					errChan <- err
				}
			}
			batch = nil
			return nil
		})
	}

	go func() {
		defer log.Printf("batcher exited ...")
		defer close(errChan)

		counter := 0

		duration := 500 * time.Millisecond
		timeout := time.NewTimer(duration)
		defer timeout.Stop()

	loop:
		for {
			timeout.Reset(duration)
			select {
			case evt, ok := <-events:
				if !ok {
					processBatch()
					break loop
				}
				batch = append(batch, evt)
				counter++
				if counter%batch_size == 0 {
					processBatch()
					log.Debugf("batch indexed %v documents", counter)
				}
			case <-timeout.C:
				log.Tracef("input channel timeout")
				if len(batch) != 0 {
					processBatch()
					log.Debugf("batch indexed %v documents", counter)
				}
			}

		}
	}()

	return errChan
}

/*
func (d DB) Add(schedulerName string, ss ...schedule.Schedule) error {
	return d.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(schedulerName))
		if err != nil {
			return fmt.Errorf("cannot create bucket %s: %s", schedulerName, err)
		}

		for _, s := range ss {
			v := b.Get([]byte(s.ID()))
			if v == nil {
				buf, err := json.Marshal([]schedule.Schedule{s})
				if err != nil {
					log.Errorf("cannot marshall schedule %v: %v", s, err)
					continue
				}
				err = b.Put([]byte(s.ID()), buf)
				if err != nil {
					log.Errorf("cannot put schedule %v: %v", s, err)
				}
				continue
			}

			buf, err := getBytes(s, v)
			if err != nil {
				log.Errorf("cannot get bytes for %v: %v", string(v), err)
				continue
			}

			err = b.Put([]byte(s.ID()), buf)
			if err != nil {
				log.Errorf("cannot put schedule %v: %v", s, err)
				continue
			}
		}

		return nil
	})
}
*/

func addToSlice(s schedule.Schedule, v []byte) ([]byte, error) {
	var arr []interface{}

	err := json.Unmarshal(v, &arr)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshall value: %v", err)
	}

	// prepend
	arr = append([]interface{}{s}, arr...)
	buf, err := json.Marshal(arr)
	if err != nil {
		return nil, fmt.Errorf("cannot Marshal slice: %v", err)
	}

	return buf, nil
}
