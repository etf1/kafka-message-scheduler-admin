package db

import (
	"github.com/etf1/kafka-message-scheduler-admin/server/sort"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler/schedule"
)

type DB interface {
	store.Store
	Search(q SearchQuery) (int, chan schedule.Schedule, error)
}

type SearchQuery struct {
	Limit
	Filter
	sort.SortBy
}

type Filter struct {
	SchedulerName string
	ScheduleID    string
	EpochRange
}
type EpochRange struct {
	From int64
	To   int64
}

type Limit struct {
	Max int
}
