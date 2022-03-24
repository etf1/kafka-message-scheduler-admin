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
	logErr           = helper.LogErr
)

func newScheduler(name string) httpresolver.Scheduler {
	localhostIP := func() (a, b, c, d byte) {
		a = 127
		b = 0
		c = 0
		d = 1
		return a, b, c, d
	}
	return httpresolver.Scheduler{
		HostName: name,
		HTTPPort: "8000",
		Instances: []httpresolver.Instance{
			{
				IP:               net.IPv4(localhostIP()),
				HostNames:        []string{"localhost"},
				Topics:           []string{"schedules"},
				HistoryTopic:     "history",
				BootstrapServers: "localhost:9092",
			},
		},
	}
}

type Runner struct {
	stopChan chan bool
	exitChan chan bool
}

func NewRunner() *Runner {
	return &Runner{
		stopChan: make(chan bool),
		exitChan: make(chan bool),
	}
}

func (r Runner) Close() {
	r.stopChan <- true
	<-r.exitChan
	log.Printf("after close return")
}

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func ARandInt() int {
	min, max := 1, 100
	return randInt(min, max)
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
			targetTopic := fmt.Sprintf("%v%v", ksch.TargetTopic(), ARandInt())
			targetID := fmt.Sprintf("%v%v", ksch.TargetKey(), ARandInt())
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
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovering from panic in runner: %v", r)
		}
		r.exitChan <- true
		log.Printf("mini runner stopped")
	}()

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
		targetTopic := fmt.Sprintf("target-topic-%v", ARandInt())
		targetKey := fmt.Sprintf("target-id-%v", i+1)
		value := "some french char: éàçèùäâ" + lipsum()
		schs[i] = newKafkaSchedule("schedules", fmt.Sprintf("schedule-%v", i+1), value, t.Unix(), targetTopic, targetKey)
	}

	logErr(liveStore.Add(sch1.Name(), genRandVersions(schs[0:10])...))
	logErr(liveStore.Add(sch2.Name(), genRandVersions(schs[100:120])...))
	logErr(liveStore.Add(sch3.Name(), genRandVersions(schs[200:230])...))

	logErr(coldStore.Add(sch1.Name(), genRandVersions(schs[0:100])...))
	logErr(coldStore.Add(sch2.Name(), genRandVersions(schs[100:200])...))
	logErr(coldStore.Add(sch3.Name(), genRandVersions(schs[200:300])...))

	logErr(historyStore.Add(sch1.Name(), genRandVersions(schs[0:50])...))
	logErr(historyStore.Add(sch2.Name(), genRandVersions(schs[100:150])...))
	logErr(historyStore.Add(sch3.Name(), genRandVersions(schs[200:250])...))

	srv := runner.NewServer(coldDB, liveDB, historyDB, resolver)

	helper.StartupHTTPServer(srv)
	<-r.stopChan
	helper.LogErr(helper.ShutdownHTTPServer(srv))

	return nil
}
