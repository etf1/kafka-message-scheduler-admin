package simple

import (
	"strings"

	stdsort "sort"

	"github.com/etf1/kafka-message-scheduler-admin/server/db"
	"github.com/etf1/kafka-message-scheduler-admin/server/sort"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler/schedule"
)

const (
	DefaultMax = 300
	ChanSize   = 1000
)

type DB struct {
	store.Store
}

func (d DB) Search(q db.SearchQuery) (total int, result chan schedule.Schedule, err error) {
	matches := func(q db.SearchQuery, sch schedule.Schedule) bool {
		match := true

		if strings.TrimSpace(q.Filter.ScheduleID) != "" {
			match = match && strings.Contains(strings.TrimSpace(sch.ID()), q.Filter.ScheduleID)
		}

		if q.EpochRange.From != 0 {
			epoch := sch.Epoch()
			match = match && q.EpochRange.From <= epoch
		}

		if q.EpochRange.To != 0 && q.EpochRange.From <= q.EpochRange.To {
			epoch := sch.Epoch()
			match = match && epoch <= q.EpochRange.To
		}
		return match
	}

	found := 0
	schedules, err := d.List(q.SchedulerName)
	if err != nil {
		return found, result, err
	}

	arr := []schedule.Schedule{}
	for sch := range schedules {
		if matches(q, sch) {
			arr = append(arr, sch)
		}
	}

	stdsort.Sort(sort.NewSort(arr, q.SortBy))

	result = make(chan schedule.Schedule, ChanSize)

	max := DefaultMax
	if q.Max > 0 {
		max = q.Max
	}

	go func() {
		defer close(result)

		for _, sch := range arr {
			if found >= max {
				break
			}
			result <- sch
			found++
		}
	}()

	return len(arr), result, nil
}
