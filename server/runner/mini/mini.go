package mini

// Kafka runner for the scheduler

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/db/simple"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/slice"
	"github.com/etf1/kafka-message-scheduler-admin/server/restapi"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/hmap"
	"github.com/etf1/kafka-message-scheduler/schedule"
	"github.com/etf1/kafka-message-scheduler/schedule/kafka"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	newKafkaSchedule = helper.NewKafkaSchedule
	lipsum           = helper.Lipsum
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
	cold := hmap.NewStore()
	live := hmap.NewStore()
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
		schs[i] = newKafkaSchedule("schedules", fmt.Sprintf("schedule-%v", i+1), lipsum(), t.Unix(), targetTopic, targetKey)
	}

	live.Add(sch1.Name(), genRandVersions(schs[0:10])...)
	live.Add(sch2.Name(), genRandVersions(schs[100:120])...)
	live.Add(sch3.Name(), genRandVersions(schs[200:230])...)

	cold.Add(sch1.Name(), genRandVersions(schs[0:100])...)
	cold.Add(sch2.Name(), genRandVersions(schs[100:200])...)
	cold.Add(sch3.Name(), genRandVersions(schs[200:300])...)

	if config.APIServerOnly() {
		srv := restapi.NewServer(simple.DB{
			Store: cold,
		}, simple.DB{
			Store: live,
		}, resolver)
		srv.Start(config.ServerAddr())
		defer srv.Stop()
	} else {
		mainRouter := mux.NewRouter().StrictSlash(true)

		mainRouter.PathPrefix("/api").Handler(http.StripPrefix("/api", restapi.NewRouter(simple.DB{
			Store: cold,
		}, simple.DB{
			Store: live,
		}, resolver)))

		mainRouter.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(config.StaticFilesDir()))))

		addr := config.ServerAddr()
		srv := &http.Server{
			Handler:      mainRouter,
			Addr:         addr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		defer helper.ShutdownHttpServer(srv, config.ShutdownTimeout())

		go func() {
			log.Printf("starting server on %s", addr)
			defer log.Printf("server stopped")
			if err := srv.ListenAndServe(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	<-r.stopChan

	log.Printf("mini server closed")

	return nil
}
