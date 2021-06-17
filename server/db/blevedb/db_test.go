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
	"github.com/etf1/kafka-message-scheduler-admin/server/store/hmap"
	"github.com/etf1/kafka-message-scheduler/schedule/simple"
)

func initDB(t *testing.T) (*hmap.Hmap, db.DB, func()) {
	sourceStore := hmap.NewStore()
	dir := helper.GenRandString("db-")
	path := dir + "/schedules.bleve"
	bdb, err := blevedb.NewDB(blevedb.Config{
		SourceStore:   sourceStore,
		InternalStore: hmap.NewStore(),
		Path:          path,
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
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
		// should default to default max
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
				t.Errorf("unexpected error: %v", err)
				return
			}

			if found != tt.expectedTotal {
				t.Errorf("unexpected total found count: %v", found)
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
				t.Errorf("unexpected result size: %v", count)
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
		sortBy      sort.By
		expectedIDs []string
	}{
		{sort.By{Field: sort.Timestamp, Order: sort.Desc}, []string{"schedule-2", "schedule-1"}},
		{sort.By{Field: sort.Timestamp, Order: sort.Asc}, []string{"schedule-1", "schedule-2"}},
		{sort.By{Field: sort.Epoch, Order: sort.Desc}, []string{"schedule-1", "schedule-2"}},
		{sort.By{Field: sort.Epoch, Order: sort.Asc}, []string{"schedule-2", "schedule-1"}},
		{sort.By{Field: sort.ID, Order: sort.Desc}, []string{"schedule-2", "schedule-1"}},
		{sort.By{Field: sort.ID, Order: sort.Asc}, []string{"schedule-1", "schedule-2"}},
		// partial
		{sort.By{Order: sort.Desc}, []string{"schedule-2", "schedule-1"}},
		{sort.By{Order: sort.Asc}, []string{"schedule-1", "schedule-2"}},
		{sort.By{Field: sort.Timestamp}, []string{"schedule-1", "schedule-2"}},
		{sort.By{Field: sort.Epoch}, []string{"schedule-2", "schedule-1"}},
		{sort.By{Field: sort.ID}, []string{"schedule-1", "schedule-2"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			_, lst, err := bdb.Search(db.SearchQuery{
				SortBy: tt.sortBy,
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			index := 0
			for s := range lst {
				if index >= len(tt.expectedIDs) {
					t.Errorf("unexpected result length > %v", len(tt.expectedIDs))
				}
				if tt.expectedIDs[index] != s.ID() {
					t.Errorf("unexpected id: %v", tt.expectedIDs[index])
				}
				index++
			}
			if index != len(tt.expectedIDs) {
				t.Errorf("unexpected result length: %v", index)
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
		{"scheduler", []string{}},
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
				t.Errorf("unexpected error: %v", err)
			}

			index := 0
			for s := range lst {
				if index >= len(tt.expectedIDs) {
					t.Errorf("unexpected result length > %v", len(tt.expectedIDs))
				}
				if tt.expectedIDs[index] != s.ID() {
					t.Errorf("unexpected id: %v", tt.expectedIDs[index])
				}
				index++
			}

			if index != len(tt.expectedIDs) {
				t.Errorf("unexpected result length: %v", index)
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
	data.Add("scheduler-2", simple.NewSchedule("video:apps:online", now, now.Add(3*time.Second)))
	data.Add("scheduler-2", simple.NewSchedule("video:apps:offline", now, now.Add(4*time.Second)))
	data.Add("scheduler-2", simple.NewSchedule("video:web:online", now, now.Add(5*time.Second)))
	data.Add("scheduler-2", simple.NewSchedule("collection:01abc-12rfg-34kl:online", now, now.Add(6*time.Second)))
	// wait for goroutines to be scheduled
	time.Sleep(1 * time.Second)

	tests := []struct {
		scheduleID    string
		schedulerName string
		expectedIDs   []string
	}{
		// no scheduler, no results
		{"", "", []string{"schedule-1", "schedule-2", "schedule-3", "video:apps:online", "video:apps:offline", "video:web:online", "collection:01abc-12rfg-34kl:online"}},
		// no search terms, all results
		{"", "scheduler-1", []string{"schedule-1"}},
		{"", "scheduler-2", []string{"schedule-2", "schedule-3", "video:apps:online", "video:apps:offline", "video:web:online", "collection:01abc-12rfg-34kl:online"}},
		{"+video", "scheduler-2", []string{"video:apps:online", "video:apps:offline", "video:web:online"}},
		{"+video +online", "scheduler-2", []string{"video:apps:online", "video:web:online"}},
		{"+video +online -web", "scheduler-2", []string{"video:apps:online"}},
		{"+video +online -web", "scheduler-2", []string{"video:apps:online"}},
		{"video:apps:online", "scheduler-2", []string{"video:apps:online"}},
		{"video*", "scheduler-2", []string{"video:apps:online", "video:apps:offline", "video:web:online"}},
		{"vid*", "scheduler-2", []string{"video:apps:online", "video:apps:offline", "video:web:online"}},
		{"vid* on*", "scheduler-2", []string{"video:apps:online", "video:web:online"}},
		{"vid* off*", "scheduler-2", []string{"video:apps:offline"}},
		{"collection:01abc-12rfg-34kl:online", "scheduler-2", []string{"collection:01abc-12rfg-34kl:online"}},
		{"collection:01abc-12rfg-34kl", "scheduler-2", []string{"collection:01abc-12rfg-34kl:online"}},
		{"coll*", "scheduler-2", []string{"collection:01abc-12rfg-34kl:online"}},
		{"collection 01abc", "scheduler-2", []string{"collection:01abc-12rfg-34kl:online"}},
		{"01abc-12rfg-34kl", "scheduler-2", []string{"collection:01abc-12rfg-34kl:online"}},
		{"01abc-12rfg", "scheduler-2", []string{"collection:01abc-12rfg-34kl:online"}},
		{"01abc", "scheduler-2", []string{"collection:01abc-12rfg-34kl:online"}},
		{"12rfg", "scheduler-2", []string{"collection:01abc-12rfg-34kl:online"}},
		{"34kl", "scheduler-2", []string{"collection:01abc-12rfg-34kl:online"}},
		{"34klx", "scheduler-2", []string{}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			total, lst, err := bdb.Search(db.SearchQuery{
				Filter: db.Filter{
					ScheduleID:    tt.scheduleID,
					SchedulerName: tt.schedulerName,
				},
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(tt.expectedIDs) == 0 && total != 0 || total < len(tt.expectedIDs) {
				t.Errorf("unexpected total: %v", total)
			}

			index := 0
			for s := range lst {
				if index >= len(tt.expectedIDs) {
					t.Errorf("unexpected result length (%v)> %v", index, len(tt.expectedIDs))
				}
				if tt.expectedIDs[index] != s.ID() {
					t.Errorf("unexpected id: %v", s.ID())
				}
				index++
			}
			if index != len(tt.expectedIDs) {
				t.Errorf("unexpected result length: %v", index)
			}
		})
	}
}

// Rule #5: search by scheduler id should filter the result list
func TestBleveDBSearch_by_epoch(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	data, bdb, clean := initDB(t)
	defer clean()

	now := time.Now()
	epoch1 := now.Add(1 * time.Second).Unix()
	epoch2 := now.Add(2 * time.Second).Unix()
	epoch3 := now.Add(3 * time.Second).Unix()
	epoch4 := now.Add(4 * time.Second).Unix()

	data.Add("scheduler-1", simple.NewSchedule("schedule-1", epoch1, now.Add(1*time.Second)))
	data.Add("scheduler-1", simple.NewSchedule("schedule-2", epoch2, now.Add(2*time.Second)))
	data.Add("scheduler-1", simple.NewSchedule("schedule-3", epoch3, now.Add(3*time.Second)))
	data.Add("scheduler-1", simple.NewSchedule("schedule-4", epoch4, now.Add(4*time.Second)))

	// wait for goroutines to be scheduled
	time.Sleep(1 * time.Second)

	tests := []struct {
		schedulerName string
		epochFrom     int64
		epochTo       int64
		expectedIDs   []string
	}{
		{"scheduler-1", 0, 0, []string{"schedule-1", "schedule-2", "schedule-3", "schedule-4"}},
		{"scheduler-1", epoch1, 0, []string{"schedule-1", "schedule-2", "schedule-3", "schedule-4"}},
		{"scheduler-1", epoch2, 0, []string{"schedule-2", "schedule-3", "schedule-4"}},
		{"scheduler-1", epoch2, epoch3, []string{"schedule-2", "schedule-3"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			srch := db.SearchQuery{
				Filter: db.Filter{
					SchedulerName: tt.schedulerName,
					EpochRange:    db.EpochRange{},
				},
			}
			if tt.epochFrom != 0 {
				srch.EpochRange.From = tt.epochFrom
			}
			if tt.epochTo != 0 {
				srch.EpochRange.To = tt.epochTo
			}
			_, lst, err := bdb.Search(srch)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			index := 0
			for s := range lst {
				if index >= len(tt.expectedIDs) {
					t.Errorf("unexpected result length > %v", len(tt.expectedIDs))
				}
				if tt.expectedIDs[index] != s.ID() {
					t.Errorf("unexpected id: %v", s.ID())
				}
				index++
			}
			if index != len(tt.expectedIDs) {
				t.Errorf("unexpected result length: %v", index)
			}
		})
	}
}

func TestBleveDBSearch_reset_store(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	data, bdb, clean := initDB(t)
	defer clean()

	bldb, ok := bdb.(blevedb.DB)
	if !ok {
		t.Errorf("unexpected type: %T", bldb)
	}

	nbSchedules := 10
	data.Add("scheduler-1", helper.SimpleSchedules(nbSchedules)...)
	data.Add("scheduler-2", helper.SimpleSchedules(nbSchedules)...)

	// wait for goroutines to be scheduled
	time.Sleep(1 * time.Second)

	tests := []struct {
		schedulerName        string
		expectedInitialCount int
	}{
		{"scheduler-1", nbSchedules},
	}

	checkCount := func(scheduler string, expectedCount int) {
		srch := db.SearchQuery{
			Filter: db.Filter{
				SchedulerName: scheduler,
			},
		}
		count, _, err := bdb.Search(srch)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if count != expectedCount {
			t.Errorf("unexpected count:%v expected:%v", count, expectedCount)
		}
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			checkCount(tt.schedulerName, tt.expectedInitialCount)
			data.Reset(tt.schedulerName)
			time.Sleep(1 * time.Second)
			checkCount(tt.schedulerName, 0)
			checkCount("scheduler-2", nbSchedules)
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
				t.Errorf("unexpected error: %v", err)
			}

			index := 0
			for s := range lst {
				if index >= len(tt.expectedIDs) {
					t.Errorf("unexpected result length > %v", len(tt.expectedIDs))
				}
				if tt.expectedIDs[index] != s.ID() {
					t.Errorf("unexpected id: %v", s.ID())
				}
				index++
			}

			t.Logf("index=%v expected_length=%v found=%v", index, len(tt.expectedIDs), found)
			if index != len(tt.expectedIDs) {
				t.Errorf("unexpected result length: %v", index)
			}
		})
	}
}
