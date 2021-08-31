package store

import (
	"fmt"

	"github.com/etf1/kafka-message-scheduler/schedule"
)

type Schedule struct {
	SchedulerName     string `json:"scheduler"`
	schedule.Schedule `json:"schedule"`
}

func (s Schedule) String() string {
	return fmt.Sprintf("{scheduler:%q schedule:%s}", s.SchedulerName, s.Schedule)
}

type EventType int

const (
	UpsertType EventType = iota
	DeletedType
	// informs that the database has been reset
	StoreResetType
)

type Event struct {
	EventType
	Schedule
}

type Watchable interface {
	Watch() (chan Event, error)
}

type Batchable interface {
	Batch(events chan Event) chan error
}

type Store interface {
	Get(schedulerName string, scheduleID string) ([]Schedule, error)
	List(schedulerName string) (chan Schedule, error)
}

type MutableStore interface {
	Store
	Add(schedulerName string, ss ...schedule.Schedule) error
	Delete(schedulerName string, ss ...schedule.Schedule) error
}

type WatchableStore interface {
	Store
	Watchable
}

type BatchableStore interface {
	Store
	Batchable
}
