package kafka

import (
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/hmap"
	"github.com/etf1/kafka-message-scheduler/schedule"
	"github.com/etf1/kafka-message-scheduler/schedule/kafka"
	log "github.com/sirupsen/logrus"
)

type event struct {
	name string
	*confluent.Message
}

type consumer struct {
	name             string
	consumer         *confluent.Consumer
	bootstrapServers string
	topics           []string
	stopChan         chan bool
}

func newConsumer(name, bootstrapServers string, topics []string) (consumer, error) {
	log.Printf("new consumer topics=%v bootstrapServers=%v", topics, bootstrapServers)
	kafkaConsumer, err := confluent.NewConsumer(&confluent.ConfigMap{
		"bootstrap.servers":  bootstrapServers,
		"group.id":           "scheduler-admin-cg",
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
		make(chan bool),
	}, nil
}

func (c consumer) processMessage(events chan event) {
	e := c.consumer.Poll(100)

	if e == nil {
		return
	}
	switch evt := e.(type) {
	case *confluent.Message:
		events <- event{
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
		defer log.Printf("consumer closed: %+v", c)
		defer c.consumer.Close()
		for {
			select {
			case <-c.stopChan:
				log.Printf("consumer closing")
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
}

type processor struct {
	processChan   chan event
	processedChan chan event
	action        func(evt event) error
	stopChan      chan bool
}

func newProcessor(action func(evt event) error) processor {
	return processor{
		processChan:   make(chan event, 10000),
		processedChan: make(chan event, 10000),
		action:        action,
		stopChan:      make(chan bool),
	}
}

func (p processor) close() {
	p.stopChan <- true
}

func (p processor) start() {
	go func() {
		defer log.Printf("processor closed")
		defer close(p.processChan)
		for {
			select {
			case <-p.stopChan:
				return
			case msg := <-p.processChan:
				log.Printf("received message: %+v", msg)
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

	s.processor.close()

	log.Printf("closing kafka store ...")
	for _, c := range s.consumers {
		c.close()
	}
}

func (s Store) Get(schedulerName string, scheduleID string) ([]store.Schedule, error) {
	return s.data.Get(schedulerName, scheduleID)
}

func (s Store) List(schedulerName string) (chan store.Schedule, error) {
	return s.data.List(schedulerName)
}

func (s Store) Add(schedulerName string, ss ...schedule.Schedule) error {
	return s.data.Add(schedulerName, ss...)
}

func (s Store) Watch() (chan store.Event, error) {
	resultChan := make(chan store.Event, 10000)

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
