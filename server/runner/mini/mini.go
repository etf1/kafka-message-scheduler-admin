package mini

// Kafka runner for the scheduler

import (
	"fmt"
	"net"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/db/simple"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/slice"
	"github.com/etf1/kafka-message-scheduler-admin/server/restapi"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/hmap"
	"github.com/etf1/kafka-message-scheduler/schedule"
	log "github.com/sirupsen/logrus"
)

var (
	newKafkaSchedule = helper.NewKafkaSchedule
)

func newScheduler(name string) httpresolver.Scheduler {
	return httpresolver.Scheduler{
		HostName: name,
		HTTPPort: "8080",
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

func (r *Runner) Start() error {
	cold := hmap.NewStore()
	live := hmap.NewStore()
	res := slice.NewResolver()

	sch1 := newScheduler("scheduler-1")
	sch2 := newScheduler("scheduler-2")
	sch3 := newScheduler("scheduler-3")

	res.Add(sch1, sch2, sch3)

	now := time.Now()

	size := 300
	schs := make([]schedule.Schedule, size)
	for i := 0; i < size; i++ {
		t := now.Add(time.Duration(i) * time.Second)
		schs[i] = newKafkaSchedule("schedules", fmt.Sprintf("schedule-%v", i+1), fmt.Sprintf("value for schedule-%v", i+1), t.Unix(), "target-topic", fmt.Sprintf("target-id-%v", i+1))
	}

	live.Add(sch1.Name(), schs[:10]...)
	live.Add(sch2.Name(), schs[100:120]...)
	live.Add(sch3.Name(), schs[200:230]...)

	cold.Add(sch1.Name(), schs[:100]...)
	cold.Add(sch2.Name(), schs[100:200]...)
	cold.Add(sch3.Name(), schs[200:300]...)

	srv := restapi.NewServer(simple.DB{
		Store: cold,
	}, simple.DB{
		Store: live,
	}, res)

	srv.Start(config.APIServerAddr())

	<-r.stopChan

	srv.Stop()

	log.Printf("mini server closed")

	return nil
}
