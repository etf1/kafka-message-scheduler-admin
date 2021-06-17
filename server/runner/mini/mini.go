package mini

// Kafka runner for the scheduler

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/db/simple"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/slice"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/hmap"
	"github.com/etf1/kafka-message-scheduler/schedule"
	"github.com/etf1/kafka-message-scheduler/schedule/kafka"
	log "github.com/sirupsen/logrus"
)

var (
	newKafkaSchedule = helper.NewKafkaSchedule
	lipsum           = helper.Lipsum
)

func newScheduler(name string) httpresolver.Scheduler {
	return httpresolver.Scheduler{
		HostName: name,
		HTTPPort: "8000",
		Instances: []httpresolver.Instance{
			{
				IP:               net.IPv4(127, 0, 0, 1),
				HostNames:        []string{"localhost"},
				Topics:           []string{"schedules"},
				BootstrapServers: "localhost:9092",
			},
		},
	}
}

type Runner struct {
	stopChan chan bool
}

func NewRunner() *Runner {
	return &Runner{
		stopChan: make(chan bool, 1),
	}
}

func (r Runner) Close() {
	r.stopChan <- true
}

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func genRandVersions(src []kafka.Schedule) []schedule.Schedule {
	min := 0
	max := 10
	result := []schedule.Schedule{}
	for i, ksch := range src {
		result = append(result, ksch)
		// to set a random deleted version
		deletedIndex := randInt(min, max)
		// number of versions
		nbVersions := randInt(min, max)
		for j := 0; j < nbVersions; j++ {
			targetTopic := fmt.Sprintf("%v%v", ksch.TargetTopic(), randInt(1, 10))
			targetID := fmt.Sprintf("%v%v", ksch.TargetKey(), randInt(1, 100))
			var value interface{} = lipsum()
			if i == deletedIndex {
				value = nil
			}
			result = append(result, newKafkaSchedule(ksch.Topic(), ksch.ID(), value, ksch.Epoch(), targetTopic, targetID))
		}
	}
	return result
}

func (r *Runner) Start() error {
	defer log.Printf("mini runner stopped")

	coldStore := hmap.NewStore()
	coldDB := simple.DB{Store: coldStore}

	liveStore := hmap.NewStore()
	liveDB := simple.DB{Store: liveStore}

	historyStore := hmap.NewStore()
	historyDB := simple.DB{Store: historyStore}

	resolver := slice.NewResolver()

	sch1 := newScheduler("scheduler-1")
	sch2 := newScheduler("scheduler-2")
	sch3 := newScheduler("scheduler-3")

	resolver.Add(sch1, sch2, sch3)

	now := time.Now()

	size := 300
	schs := make([]kafka.Schedule, size)

	for i := 0; i < size; i++ {
		t := now.Add(time.Duration(i) * time.Second)
		targetTopic := fmt.Sprintf("target-topic-%v", randInt(0, 100))
		targetKey := fmt.Sprintf("target-id-%v", i+1)
		schs[i] = newKafkaSchedule("schedules", fmt.Sprintf("schedule-%v", i+1), "some french char: éàçèùäâ"+lipsum(), t.Unix(), targetTopic, targetKey)
	}

	liveStore.Add(sch1.Name(), genRandVersions(schs[0:10])...)
	liveStore.Add(sch2.Name(), genRandVersions(schs[100:120])...)
	liveStore.Add(sch3.Name(), genRandVersions(schs[200:230])...)

	coldStore.Add(sch1.Name(), genRandVersions(schs[0:100])...)
	coldStore.Add(sch2.Name(), genRandVersions(schs[100:200])...)
	coldStore.Add(sch3.Name(), genRandVersions(schs[200:300])...)

	historyStore.Add(sch1.Name(), genRandVersions(schs[0:50])...)
	historyStore.Add(sch2.Name(), genRandVersions(schs[100:150])...)
	historyStore.Add(sch3.Name(), genRandVersions(schs[200:250])...)

	srv := runner.NewServer(coldDB, liveDB, historyDB, resolver)

loop:
	select {
	case <-helper.StartupHttpServer(srv):
		log.Printf("exitChan received")
		break loop
	case <-r.stopChan:
		log.Printf("stopChan received")
		helper.ShutdownHttpServer(srv)
	}

	return nil
}
