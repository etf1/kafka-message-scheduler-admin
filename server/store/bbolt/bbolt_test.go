package bbolt_test

import (
	"os"
	"testing"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/bbolt"
	simple_schedule "github.com/etf1/kafka-message-scheduler/schedule/simple"
)

var (
	simpleSchedule = simple_schedule.NewSchedule
)

func TestBboltStore_add_get_list(t *testing.T) {
	file := helper.GenRandString("db-")
	defer func() {
		err := os.Remove(file)
		if err != nil {
			t.Errorf("unable to delete db file %v: %v", file, err)
		}
	}()

	db, err := bbolt.NewStore(file)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer db.Close()

	now := time.Now()
	sch1 := simpleSchedule("schedule-1", now)
	sch1Bis := simpleSchedule("schedule-1", now.Add(1*time.Hour))

	sch2 := simpleSchedule("schedule-2", now)
	sch3 := simpleSchedule("schedule-3", now)
	sch4 := simpleSchedule("schedule-4", now)

	err = db.Add("scheduler-1", sch1, sch2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	err = db.Add("scheduler-2", sch3, sch4)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	err = db.Add("scheduler-3", sch1, sch1Bis)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	db.Add("scheduler-4", sch1, sch2, sch3, sch4)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	lst, err := db.Get("scheduler-1", "schedule-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(lst) != 1 {
		t.Errorf("unexpected result length: %v", len(lst))
	}

	if lst[0].ID() != "schedule-1" || lst[0].Epoch() == 0 || lst[0].Timestamp() == 0 {
		t.Errorf("unexpected schedule: %v", lst[0])
	}

	lst, err = db.Get("scheduler-3", "schedule-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(lst) != 2 {
		t.Errorf("unexpected result length: %v", len(lst))
	}

	// first schedule should be the last added, so sch1Bis
	if lst[0].ID() != "schedule-1" || lst[0].Epoch() == 0 || lst[0].Timestamp() == 0 || lst[0].Epoch() != sch1Bis.Epoch() {
		t.Errorf("unexpected schedule: %v", lst[0])
	}

	if lst[1].ID() != "schedule-1" || lst[1].Epoch() == 0 || lst[1].Timestamp() == 0 {
		t.Errorf("unexpected schedule: %v", lst[1])
	}

	schan, err := db.List("scheduler-4")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	count := 0
	for {
		_, ok := <-schan
		if !ok {
			break
		}
		count++
	}

	if count != 4 {
		t.Errorf("unexpected count: %v", count)
	}
}

func TestBboltStore_delete(t *testing.T) {
	file := helper.GenRandString("db-")
	defer func() {
		err := os.Remove(file)
		if err != nil {
			t.Errorf("unable to delete db file %v: %v", file, err)
		}
	}()

	db, err := bbolt.NewStore(file)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer db.Close()

	now := time.Now()
	sch1 := simpleSchedule("schedule-1", now)

	db.Add("scheduler-1", sch1)

	lst, err := db.Get("scheduler-1", "schedule-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(lst) != 1 {
		t.Errorf("unexpected result: %v", len(lst))
	}

	err = db.Delete("scheduler-1", sch1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	lst, err = db.Get("scheduler-1", "schedule-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// should be 0 because deleted
	if len(lst) != 0 {
		t.Errorf("unexpected result: %v", len(lst))
	}

	schan, err := db.List("scheduler-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	count := 0
	for {
		_, ok := <-schan
		if !ok {
			break
		}
		count++
	}

	if count != 0 {
		t.Errorf("unexpected count: %v", count)
	}
}

func TestBboltStore_batch(t *testing.T) {
	file := helper.GenRandString("db-")
	defer func() {
		err := os.Remove(file)
		if err != nil {
			t.Errorf("unable to delete db file %v: %v", file, err)
		}
	}()

	db, err := bbolt.NewStore(file)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer db.Close()

	now := time.Now()
	sch1 := simpleSchedule("schedule-1", now)
	sch1Bis := simpleSchedule("schedule-1", now.Add(1*time.Hour))

	sch2 := simpleSchedule("schedule-2", now)
	sch3 := simpleSchedule("schedule-3", now)
	sch4 := simpleSchedule("schedule-4", now)

	events := make(chan store.Event)

	go func() {
		defer close(events)

		events <- store.Event{
			EventType: store.UpsertType,
			Schedule: store.Schedule{
				SchedulerName: "scheduler-1",
				Schedule:      sch1,
			},
		}
		events <- store.Event{
			EventType: store.UpsertType,
			Schedule: store.Schedule{
				SchedulerName: "scheduler-1",
				Schedule:      sch1Bis,
			},
		}
		events <- store.Event{
			EventType: store.UpsertType,
			Schedule: store.Schedule{
				SchedulerName: "scheduler-2",
				Schedule:      sch2,
			},
		}
		events <- store.Event{
			EventType: store.UpsertType,
			Schedule: store.Schedule{
				SchedulerName: "scheduler-3",
				Schedule:      sch3,
			},
		}
		events <- store.Event{
			EventType: store.UpsertType,
			Schedule: store.Schedule{
				SchedulerName: "scheduler-4",
				Schedule:      sch4,
			},
		}
		events <- store.Event{
			EventType: store.DeletedType,
			Schedule: store.Schedule{
				SchedulerName: "scheduler-4",
				Schedule:      sch4,
			},
		}
	}()

	counter := 0
	errs := db.Batch(events)
loop:
	for {
		select {
		case e, ok := <-errs:
			if !ok {
				break loop
			}
			if e != nil {
				counter++
			}
		case <-time.After(2 * time.Second):
			counter = -1
			break loop
		}
	}

	// check channel is closed and no errors
	if counter == -1 {
		t.Errorf("unexpected error: channel not closed")
	}
	if counter != 0 {
		t.Errorf("unexpected length: %v", counter)
	}

	// check schedule-1 versions are present
	lst, err := db.Get("scheduler-1", "schedule-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(lst) != 2 {
		t.Errorf("unexpected result: %v", len(lst))
	}

	// check another scheduler
	lst, err = db.Get("scheduler-2", "schedule-2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(lst) != 1 {
		t.Errorf("unexpected result: %v", len(lst))
	}

	// check deleted schedule is deleted
	lst, err = db.Get("scheduler-4", "schedule-4")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(lst) != 0 {
		t.Errorf("unexpected result: %v", len(lst))
	}
}
