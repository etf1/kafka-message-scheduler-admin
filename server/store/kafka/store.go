package kafka

import (
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/hmap"
	"github.com/etf1/kafka-message-scheduler/schedule"
	"github.com/etf1/kafka-message-scheduler/schedule/kafka"
	log "github.com/sirupsen/logrus"
)

var (
	PolltimeoutMs = 100
	ChanSize      = 10000
)

type eventType int

const (
	messageType eventType = iota
	storeResetType
)

type event struct {
	evtType eventType
	name    string
	*confluent.Message
}

type consumer struct {
	name             string
	consumer         *confluent.Consumer
	bootstrapServers string
	topics           []string
	stopChan         chan bool
	exitChan         chan bool
}

func newConsumer(name, bootstrapServers string, topics []string) (consumer, error) {
	log.Printf("new consumer topics=%v bootstrapServers=%v", topics, bootstrapServers)
	kafkaConsumer, err := confluent.NewConsumer(&confluent.ConfigMap{
		"bootstrap.servers":  bootstrapServers,
		"group.id":           helper.GenRandString("kafka-store-"),
		"session.timeout.ms": 6000,
		"enable.auto.commit": false,
		"auto.offset.reset":  "earliest",
	})
	if err != nil {
		return consumer{}, err
	}

	return consumer{
		name,
		kafkaConsumer,
		bootstrapServers,
		topics,
		make(chan bool, 1),
		make(chan bool, 1),
	}, nil
}

func (c consumer) processMessage(events chan event) {
	e := c.consumer.Poll(PolltimeoutMs)

	if e == nil {
		return
	}
	switch evt := e.(type) {
	case *confluent.Message:
		events <- event{
			messageType,
			c.name,
			evt,
		}
	case confluent.Error:
		log.Errorf("received kakfa error: %v", evt)
	default:
		log.Printf("ignored: %+v", e)
	}
}

func (c consumer) start(events chan event) error {
	err := c.consumer.SubscribeTopics(c.topics, nil)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			c.consumer.Close()
			c.exitChan <- true
			log.Printf("consumer closed: %+v", c)
		}()

		for {
			select {
			case <-c.stopChan:
				log.Printf("closing consumer")
				return
			default:
				c.processMessage(events)
			}
		}
	}()

	return nil
}

func (c consumer) close() {
	c.stopChan <- true
	<-c.exitChan
}

type processor struct {
	processChan   chan event
	processedChan chan event
	action        func(evt event) error
	stopChan      chan bool
	exitChan      chan bool
}

func newProcessor(action func(evt event) error) processor {
	return processor{
		processChan:   make(chan event, ChanSize),
		processedChan: make(chan event, ChanSize),
		action:        action,
		stopChan:      make(chan bool),
		exitChan:      make(chan bool),
	}
}

func (p processor) close() {
	p.stopChan <- true
	<-p.exitChan
	log.Printf("after processor close")
}

func (p processor) start() {
	go func() {
		defer func() {
			//close(p.processChan)
			//close(p.processedChan)
			p.exitChan <- true
			log.Printf("processor closed")
		}()

		for {
			select {
			case <-p.stopChan:
				log.Printf("processor stopChan called")
				close(p.processChan)
			case msg, ok := <-p.processChan:
				if !ok {
					close(p.processedChan)
					return
				}
				if p.action != nil {
					err := p.action(msg)
					if err != nil {
						log.Errorf("cannot process event %v : %v", msg, err)
						break
					}
				}
				p.processedChan <- msg
			}
		}
	}()
}

type Bucket struct {
	Name             string
	BootstrapServers string
	Topics           []string
}

type Store struct {
	consumers map[string]consumer
	data      store.MutableStore
	processor
}

func NewStore(buckets []Bucket) (Store, error) {
	ms := hmap.NewStore()

	action := func(evt event) error {
		return ms.Add(evt.name, kafka.Schedule{
			Message: evt.Message,
		})
	}

	p := newProcessor(action)
	p.start()

	s := Store{
		consumers: make(map[string]consumer),
		data:      ms,
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

		s.consumers[bucket.Name] = c
	}

	return s, nil
}

func (s Store) Close() {
	defer log.Printf("kafka store closed")

	log.Printf("closing kafka store ...")
	for _, c := range s.consumers {
		c.close()
	}

	// wait for consumer Poll timeout, otherwise we will get "panic: send on closed channel"
	time.Sleep(1 * time.Second)

	// s.processor.close()
}

func (s Store) Get(schedulerName, scheduleID string) ([]store.Schedule, error) {
	return s.data.Get(schedulerName, scheduleID)
}

func (s Store) List(schedulerName string) (chan store.Schedule, error) {
	return s.data.List(schedulerName)
}

func (s Store) Add(schedulerName string, ss ...schedule.Schedule) error {
	return s.data.Add(schedulerName, ss...)
}

func (s Store) Watch() (chan store.Event, error) {
	resultChan := make(chan store.Event, ChanSize)

	go func() {
		for e := range s.processedChan {
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
