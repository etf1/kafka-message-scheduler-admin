package slice

import (
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers"
)

// Slice implements Schedulers interface
type Slice struct {
	data []schedulers.Scheduler
}

type Scheduler struct {
	SchedulerName string `json:"name"`
}

func (s Scheduler) Name() string {
	return s.SchedulerName
}

func NewResolver() *Slice {
	return &Slice{
		data: make([]schedulers.Scheduler, 0),
	}
}

func (s Slice) List() ([]schedulers.Scheduler, error) {
	return append([]schedulers.Scheduler{}, s.data...), nil
}

func (s *Slice) Reset() {
	s.data = make([]schedulers.Scheduler, 0)
}

func (s *Slice) Add(sch ...schedulers.Scheduler) {
	s.data = append(s.data, sch...)
}
