package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/etf1/kafka-message-scheduler-admin/server/decoder"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler/schedule"
)

var (
	DefaultTimeout = 5 * time.Second
)

type Schedule struct {
	ScheduleID         string `json:"id"`
	ScheduleEpoch      int64  `json:"epoch"`
	ScheduleTimestamp  int64  `json:"timestamp"`
	MessageTargetTopic string `json:"target-topic"`
	MessageTargetKey   string `json:"target-key"`
	MessageTopic       string `json:"topic"`
	MessageValue       []byte `json:"value"`
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
	return fmt.Sprintf("{id:%v epoch:%v timestamp:%v}", s.ID(), s.Epoch(), s.Timestamp())
}

type HTTPRetriever struct {
	httpresolver.Resolver
	dec decoder.Decoder
}

func NewStore(r httpresolver.Resolver, dec decoder.Decoder) *HTTPRetriever {
	return &HTTPRetriever{
		Resolver: r,
		dec:      dec,
	}
}

func (h HTTPRetriever) Get(schedulerName, scheduleID string) ([]store.Schedule, error) {
	return h.getSchedules(schedulerName, func(s schedule.Schedule) bool {
		return s.ID() == scheduleID
	})
}

func (h HTTPRetriever) List(schedulerName string) (chan store.Schedule, error) {
	schedules, err := h.getSchedules(schedulerName, func(s schedule.Schedule) bool {
		return true
	})
	if err != nil {
		return nil, err
	}

	result := make(chan store.Schedule)
	go func() {
		defer close(result)

		for _, sch := range schedules {
			result <- sch
		}
	}()

	return result, nil
}

func (h HTTPRetriever) getSchedules(schedulerName string, filter func(s schedule.Schedule) bool) ([]store.Schedule, error) {
	schedulers, err := h.Resolver.List()
	if err != nil {
		return nil, err
	}

	log.Printf("schedulers=%+v", schedulers)

	result := []store.Schedule{}

	for _, _scheduler := range schedulers {
		sch, ok := _scheduler.(httpresolver.Scheduler)
		if !ok {
			log.Errorf("failed to assert type %v: %T", _scheduler, _scheduler)
			continue
		}

		log.Printf("sch.HostName=%v schedulerName=%v", sch.HostName, schedulerName)

		if sch.HostName == schedulerName {
			for _, instance := range sch.Instances {
				log.Printf("instance=%v", instance.Name()+":"+sch.HTTPPort)
				resp, err := helper.Get(instance.Name()+":"+sch.HTTPPort, "/schedules", DefaultTimeout)
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					return nil, fmt.Errorf("http request failed for %v with unexpected status code  %v", instance.Name(), resp.StatusCode)
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return nil, err
				}

				res := []Schedule{}
				err = json.NewDecoder(bytes.NewReader(body)).Decode(&res)
				if err != nil {
					return nil, err
				}

				log.Printf("schedules=%v", res)

				for _, s := range res {
					if filter(s) {
						var sch schedule.Schedule = s

						if h.dec != nil {
							sdec, err := h.dec.Decode(s)
							if err != nil {
								log.Warnf("cannot decode rest schedule %v: %v", err, s.ID())
							} else {
								sch = sdec
							}
						}

						result = append(result, store.Schedule{
							SchedulerName: schedulerName,
							Schedule:      sch,
						})
					}
				}
			}
		}
	}

	return result, nil
}
