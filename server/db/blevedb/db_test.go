// INTEGRATION TESTS

package blevedb_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/db"
	"github.com/etf1/kafka-message-scheduler-admin/server/db/blevedb"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/sort"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/hmap"
	"github.com/etf1/kafka-message-scheduler/schedule/simple"
)

func initDB(t *testing.T) (store.MutableStore, db.DB, func()) {
	sourceStore := hmap.NewStore()
	dir := helper.GenRandString("db-")
	path := dir + "/schedules.bleve"
	bdb, err := blevedb.NewDB(blevedb.Config{
		SourceStore:   sourceStore,
		InternalStore: hmap.NewStore(),
		Path:          path,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("data path created: %v", path)

	deferFunc := func() {
		bdb.Close()

		err := os.RemoveAll(dir)
		if err != nil {
			log.Fatal(err)
		}
		t.Logf("data path deleted: %v", path)
	}

	return sourceStore, bdb, deferFunc
}

// Rule #1: the max parameter should limit the result list
func TestBleveDBSearch_max(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	data, bdb, clean := initDB(t)
	defer clean()

	now := time.Now()
	data.Add("scheduler-1", simple.NewSchedule("schedule-1", now, now))
	data.Add("scheduler-2", simple.NewSchedule("schedule-2", now.Add(1*time.Second), now.Add(1*time.Second)))

	// wait for goroutines to be scheduled
	time.Sleep(1 * time.Second)

	tests := []struct {
		max                int
		expectedTotal      int
		expectedResultSize int
	}{
		{1, 2, 1},
		{2, 2, 2},
		// should default to 100
		{0, 2, 2},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			found, lst, err := bdb.Search(db.SearchQuery{
				Limit: db.Limit{
					Max: tt.max,
				},
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if found != tt.expectedTotal {
				t.Fatalf("unexpected total found count: %v", found)
			}

			count := 0
			if lst != nil {
				for {
					_, ok := <-lst
					if !ok {
						break
					}
					count++
				}
			}

			if count != tt.expectedResultSize {
				t.Fatalf("unexpected result size: %v", count)
			}
		})
	}
}

// Rule #2: the sort parameter should apply a sort on the result list
func TestBleveDBSearch_sort(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	data, bdb, clean := initDB(t)
	defer clean()

	now := time.Now()
	data.Add("scheduler-1", simple.NewSchedule("schedule-1", now.Add(1*time.Second), now))
	data.Add("scheduler-2", simple.NewSchedule("schedule-2", now, now.Add(1*time.Second)))

	// wait for goroutines to be scheduled
	time.Sleep(1 * time.Second)

	tests := []struct {
		sortBy      sort.SortBy
		expectedIDs []string
	}{
		{sort.SortBy{SortField: sort.Timestamp, SortOrder: sort.Desc}, []string{"schedule-2", "schedule-1"}},
		{sort.SortBy{SortField: sort.Timestamp, SortOrder: sort.Asc}, []string{"schedule-1", "schedule-2"}},
		{sort.SortBy{SortField: sort.Epoch, SortOrder: sort.Desc}, []string{"schedule-1", "schedule-2"}},
		{sort.SortBy{SortField: sort.Epoch, SortOrder: sort.Asc}, []string{"schedule-2", "schedule-1"}},
		{sort.SortBy{SortField: sort.ID, SortOrder: sort.Desc}, []string{"schedule-2", "schedule-1"}},
		{sort.SortBy{SortField: sort.ID, SortOrder: sort.Asc}, []string{"schedule-1", "schedule-2"}},
		// partial
		{sort.SortBy{SortOrder: sort.Desc}, []string{"schedule-2", "schedule-1"}},
		{sort.SortBy{SortOrder: sort.Asc}, []string{"schedule-1", "schedule-2"}},
		{sort.SortBy{SortField: sort.Timestamp}, []string{"schedule-1", "schedule-2"}},
		{sort.SortBy{SortField: sort.Epoch}, []string{"schedule-2", "schedule-1"}},
		{sort.SortBy{SortField: sort.ID}, []string{"schedule-1", "schedule-2"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			_, lst, err := bdb.Search(db.SearchQuery{
				SortBy: tt.sortBy,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			index := 0
			for s := range lst {
				if index >= len(tt.expectedIDs) {
					t.Fatalf("unexpected result length > %v", len(tt.expectedIDs))
				}
				if tt.expectedIDs[index] != s.ID() {
					t.Fatalf("unexpected id: %v", tt.expectedIDs[index])
				}
				index++
			}
			if index != len(tt.expectedIDs) {
				t.Fatalf("unexpected result length: %v", index)
			}
		})
	}
}

// Rule #3: search by scheduler name should filter the result list
func TestBleveDBSearch_by_scheduler_name(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	data, bdb, clean := initDB(t)
	defer clean()

	now := time.Now()
	data.Add("scheduler-1", simple.NewSchedule("schedule-1", now.Add(1*time.Second), now))
	data.Add("scheduler-2", simple.NewSchedule("schedule-2", now, now.Add(1*time.Second)))

	// wait for goroutines to be scheduled
	time.Sleep(1 * time.Second)

	tests := []struct {
		schedulerName string
		expectedIDs   []string
	}{
		{"", []string{"schedule-1", "schedule-2"}},
		{"scheduler-1", []string{"schedule-1"}},
		{"scheduler", []string{"schedule-1", "schedule-2"}},
		{"scheduler*", []string{"schedule-1", "schedule-2"}},
		{"xxx", []string{}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			_, lst, err := bdb.Search(db.SearchQuery{
				Filter: db.Filter{
					SchedulerName: tt.schedulerName,
				},
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			index := 0
			for s := range lst {
				if index >= len(tt.expectedIDs) {
					t.Fatalf("unexpected result length > %v", len(tt.expectedIDs))
				}
				if tt.expectedIDs[index] != s.ID() {
					t.Fatalf("unexpected id: %v", tt.expectedIDs[index])
				}
				index++
			}
			//if index == 0 && len(tt.expectedIDs) > 0 {
			if index != len(tt.expectedIDs) {
				t.Fatalf("unexpected result length: %v", index)
			}
		})
	}
}

// Rule #4: search by scheduler id should filter the result list
func TestBleveDBSearch_by_scheduler_id(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	data, bdb, clean := initDB(t)
	defer clean()

	now := time.Now()
	data.Add("scheduler-1", simple.NewSchedule("schedule-1", now.Add(2*time.Second), now))
	data.Add("scheduler-2", simple.NewSchedule("schedule-2", now.Add(1*time.Second), now.Add(1*time.Second)))
	data.Add("scheduler-2", simple.NewSchedule("schedule-3", now, now.Add(2*time.Second)))

	// wait for goroutines to be scheduled
	time.Sleep(1 * time.Second)

	tests := []struct {
		scheduleID    string
		schedulerName string
		expectedIDs   []string
	}{
		{"", "", []string{"schedule-1", "schedule-2", "schedule-3"}},
		{"", "scheduler-1", []string{"schedule-1"}},
		{"", "scheduler-2", []string{"schedule-2", "schedule-3"}},
		{"xxx", "scheduler-2", []string{}},
		{"schedule*", "scheduler-2", []string{"schedule-2", "schedule-3"}},
		{"schedule", "scheduler-2", []string{"schedule-2", "schedule-3"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			_, lst, err := bdb.Search(db.SearchQuery{
				Filter: db.Filter{
					ScheduleID:    tt.scheduleID,
					SchedulerName: tt.schedulerName,
				},
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			index := 0
			for s := range lst {
				if index >= len(tt.expectedIDs) {
					t.Fatalf("unexpected result length > %v", len(tt.expectedIDs))
				}
				if tt.expectedIDs[index] != s.ID() {
					t.Fatalf("unexpected id: %v", s.ID())
				}
				index++
			}
			if index != len(tt.expectedIDs) {
				t.Fatalf("unexpected result length: %v", index)
			}
		})
	}
}

// Test for issue corection when many delete events are triggered followed by an add event
// Before fix, the schedule was not stored in the internal store

func TestBleveDBSearch_delete(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	data, bdb, clean := initDB(t)
	defer clean()

	now := time.Now()

	data.Add("scheduler-1", simple.NewSchedule("schedule-1", now.Add(2*time.Second), now))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Delete("scheduler-1", simple.NewSchedule("schedule-1", now, now.Add(2*time.Second)))
	data.Add("scheduler-1", simple.NewSchedule("schedule-1", now.Add(2*time.Second), now))

	// wait for goroutines to be scheduled
	time.Sleep(1 * time.Second)

	tests := []struct {
		scheduleID    string
		schedulerName string
		expectedIDs   []string
	}{
		{"schedule-1", "scheduler-1", []string{"schedule-1"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			found, lst, err := bdb.Search(db.SearchQuery{
				Filter: db.Filter{
					ScheduleID:    tt.scheduleID,
					SchedulerName: tt.schedulerName,
				},
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			index := 0
			for s := range lst {
				if index >= len(tt.expectedIDs) {
					t.Fatalf("unexpected result length > %v", len(tt.expectedIDs))
				}
				if tt.expectedIDs[index] != s.ID() {
					t.Fatalf("unexpected id: %v", s.ID())
				}
				index++
			}

			t.Logf(">>>> index=%v expected_length=%v found=%v", index, len(tt.expectedIDs), found)
			//if index == 0 && len(tt.expectedIDs) > 0 || index != 0 && len(tt.expectedIDs) == 0 {
			if index != len(tt.expectedIDs) {
				t.Fatalf("unexpected result length: %v", index)
			}
		})
	}
}
