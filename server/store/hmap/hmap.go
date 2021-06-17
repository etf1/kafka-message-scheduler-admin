package hmap

import (
	"sync"

	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler/schedule"
)

type SchedulerName string
type ScheduleID string
type SchedulesMap map[ScheduleID][]schedule.Schedule
type SchedulerSchedules map[SchedulerName]SchedulesMap

type Hmap struct {
	mutex     *sync.RWMutex
	data      SchedulerSchedules
	watchChan chan store.Event
}

func NewStore() *Hmap {
	return &Hmap{
		data:      SchedulerSchedules{},
		mutex:     &sync.RWMutex{},
		watchChan: make(chan store.Event, 1000000),
	}
}

func (h *Hmap) Reset() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.data = SchedulerSchedules{}
}

func NewSchedulerSchedules(schedulerName string, schedules []schedule.Schedule) SchedulerSchedules {
	result := SchedulerSchedules{}
	for _, s := range schedules {
		var arr = result[SchedulerName(schedulerName)][ScheduleID(s.ID())]
		result[SchedulerName(schedulerName)][ScheduleID(s.ID())] = append(arr, s)
	}
	return result
}

func (h *Hmap) Add(schedulerName string, ss ...schedule.Schedule) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.data[SchedulerName(schedulerName)] == nil {
		h.data[SchedulerName(schedulerName)] = SchedulesMap{}
	}

	for _, sch := range ss {
		var arr = h.data[SchedulerName(schedulerName)][ScheduleID(sch.ID())]
		h.data[SchedulerName(schedulerName)][ScheduleID(sch.ID())] = append(arr, sch)
		h.watchChan <- store.Event{
			EventType: store.UpsertType,
			Schedule: store.Schedule{
				SchedulerName: schedulerName,
				Schedule:      sch,
			},
		}
	}

	return nil
}

func (h Hmap) Get(schedulerName string, scheduleID string) ([]store.Schedule, error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	lst := h.data[SchedulerName(schedulerName)][ScheduleID(scheduleID)]
	result := make([]store.Schedule, 0)
	for _, sch := range lst {
		result = append(result, store.Schedule{
			SchedulerName: schedulerName,
			Schedule:      sch,
		})
	}
	return result, nil
}

func (h Hmap) Delete(schedulerName string, ss ...schedule.Schedule) error {
	for _, sch := range ss {
		delete(h.data[SchedulerName(schedulerName)], ScheduleID(sch.ID()))
		h.watchChan <- store.Event{
			EventType: store.DeletedType,
			Schedule: store.Schedule{
				SchedulerName: schedulerName,
				Schedule:      sch,
			},
		}
	}

	if len(h.data[SchedulerName(schedulerName)]) == 0 {
		delete(h.data, SchedulerName(schedulerName))
	}

	return nil
}

func (h Hmap) List(schedulerName string) (chan store.Schedule, error) {
	h.mutex.RLock()

	result := make(chan store.Schedule, 1000)

	go func() {
		defer h.mutex.RUnlock()
		defer close(result)

		for _, schedules := range h.data[SchedulerName(schedulerName)] {
			if len(schedules) > 0 {
				result <- store.Schedule{
					SchedulerName: schedulerName,
					// in list we want to return only the "latest" version of the schedule
					Schedule: schedules[len(schedules)-1],
				}
			}
		}
	}()

	return result, nil
}

func (h Hmap) Watch() (chan store.Event, error) {
	return h.watchChan, nil
}
