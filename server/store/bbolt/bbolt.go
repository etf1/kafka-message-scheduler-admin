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
	ErrUnknowBucket = fmt.Errorf("unknown bucket")
	batch_size      = 1000
)

/*
type indexEpochItem struct {
	ScheduleId string
	Epoch      int64
}
*/
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
	/*
		indexEpoch := make([]indexEpochItem, 0)
		now := time.Now()
		for i := 0; i < 1000000; i++ {
			t := now.Add(time.Duration(i) * time.Second)
			indexEpoch = insertIndexEpoch(indexEpoch, indexEpochItem{
				ScheduleId: fmt.Sprintf("schedule-%d", i),
				Epoch:      t.Unix(),
			})
		}
	*/
	return DB{
		bbolt_db: db,
		//indexEpoch: indexEpoch,
	}, nil
}

func (d DB) Get(schedulerName string, scheduleID string) ([]store.Schedule, error) {
	var schedules []store.Schedule

	err := d.bbolt_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(schedulerName))
		if b == nil {
			// return fmt.Errorf("%w: %v", ErrUnknowBucket, schedulerName)
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

	/*
		result := make([]schedule.Schedule, len(schedules))
		for i := 0; i < len(schedules); i++ {
			//fmt.Printf(">>> %v\n", schedules[i])
			result[i] = schedules[i]
		}

		return result, err
	*/

	return schedules, nil
}

/*
func (d DB) List(schedulerName string) (chan schedule.Schedule, error) {
	result := make(chan schedule.Schedule, 1000)

	go func() {
		defer close(result)

		d.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(schedulerName))
			if b == nil {
				return fmt.Errorf("unknown bucket: %v", schedulerName)
			}

			for i := len(d.indexEpoch) - 1; i >= 0; i-- {
				schs, err := d.Get(schedulerName, d.indexEpoch[i].ScheduleId)
				if err != nil {
					return err
				}
				if len(schs) > 0 {
					result <- schs[0]
				}
			}
			return nil
		})
	}()

	return result, nil
}
*/

func (d DB) List(schedulerName string) (chan store.Schedule, error) {
	result := make(chan store.Schedule, 1000)

	go func() {
		defer close(result)

		d.bbolt_db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(schedulerName))
			if b == nil {
				// return fmt.Errorf("%w: %v", ErrUnknowBucket, schedulerName)
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

/*
func insertIndexEpoch(data []indexEpochItem, el indexEpochItem) []indexEpochItem {
	index := sort.Search(len(data), func(i int) bool { return data[i].Epoch > el.Epoch })
	data = append(data, indexEpochItem{})
	copy(data[index+1:], data[index:])
	data[index] = el
	return data
}
*/

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
				log.Printf("bbolt store put (1) %s size %v", s.ID(), len(buf))
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
			log.Printf("bbolt store put (2) %s size %v", s.ID(), len(buf))
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

/*
func (d DB) removeFromBucket(bucketName string, sch schedule.Schedule) error {
	return d.bbolt_db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}
		err := b.Delete([]byte(sch.ID()))
		if err != nil {
			log.Errorf("cannot delete schedule %v in bucket %v: %v", sch, b, err)
		}

		return nil
	})
}

func (d DB) addToBucket(bucketName string, sch schedule.Schedule) error {
	return d.bbolt_db.Batch(func(tx *bolt.Tx) error {
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
	})
}

func (d DB) Batch(event chan store.Event) chan error {
	result := make(chan error, batch_size)

	go func() {
		defer log.Warnf("batcher exited ...")
		defer close(result)
		counter := 0
		for evt := range event {
			//log.Warnf("batcher received event: %v", evt)

			var err error
			switch evt.EventType {
			case store.DeletedType:
				err = d.removeFromBucket(evt.SchedulerName, evt.Schedule.Schedule)
			case store.UpsertType:
				err = d.addToBucket(evt.SchedulerName, evt.Schedule.Schedule)
			}
			if err != nil {
				result <- err
			}
			counter++
			if counter%batch_size == 0 {
				log.Warnf("batch processed %v documents", counter)
			}
		}
	}()

	return result
}
*/

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
		defer log.Warnf("batcher exited ...")
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
					log.Warnf("batch indexed %v documents", counter)
				}
			case <-timeout.C:
				log.Warnf("input channel timeout")
				if len(batch) != 0 {
					processBatch()
					log.Warnf("batch indexed %v documents", counter)
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
