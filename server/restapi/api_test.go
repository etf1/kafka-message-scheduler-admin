// // FUNCTIONAL TESTS
package restapi_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/slice"
	"github.com/etf1/kafka-message-scheduler-admin/server/restapi"
	"github.com/etf1/kafka-message-scheduler/schedule"
)

// Rule #1: get schedule response when not found: 404, 500
func TestRestAPIServer_getSchedule_not_found_or_error(t *testing.T) {
	router, stores, _ := newRouter()

	s1 := simpleSchedule("schedule-1", time.Now().Unix())

	scheduler1 := schedulerSchedules{
		"scheduler-1",
		schedulesSlice(s1),
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	tests := []struct {
		schedules        []schedulerSchedules
		url              string
		expectedCode     int
		expectedResponse string
	}{
		// no data in store
		{nil, "/scheduler/scheduler-0/schedule/schedule-0", http.StatusNotFound, ""},
		// some data in store
		{schedulersSchedules(scheduler1), "/scheduler/scheduler-0/schedule/schedule-0", http.StatusNotFound, ""},
		// no data in store (existing scheduler)
		{nil, fmt.Sprintf("/scheduler/%s/schedule/schedule-0", scheduler1.SchedulerName), http.StatusNotFound, ""},
		// some data in store (existing scheduler)
		{schedulersSchedules(scheduler1), fmt.Sprintf("/scheduler/%s/schedule/schedule-0", scheduler1.SchedulerName), http.StatusNotFound, ""},
		// TODO: with error
		// {[]schedule.Schedule{s1}, fmt.Sprintf("/scheduler/%s/schedule/schedule-0", "ERROR"), http.StatusInternalServerError, `{"error":"simulated error"}`},
	}

	for _, prefix := range []string{"", "/live"} {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
				createSchedulerSchedules(tt.schedules, stores...)

				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, prefix+tt.url, http.NoBody)
				response := executeRequest(router, req)
				checkResponseJSON(t, tt.expectedCode, response, tt.expectedResponse)
			})
		}
	}
}

// Rule #2: get schedule response when found: 302 + expected payload
func TestRestAPIServer_getSchedule_found(t *testing.T) {
	router, stores, _ := newRouter()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	s1 := newSchedule("scheduler-1", "schedule-1", time.Now().Unix())
	s2 := newSchedule("scheduler-1", "schedule-2", time.Now().Unix())

	scheduler1 := schedulerSchedules{
		"scheduler-1",
		schedulesSlice(s1, s2),
	}

	scheduler1bis := schedulerSchedules{
		"scheduler-1",
		schedulesSlice(s1, s1, s2),
	}

	url := fmt.Sprintf("/scheduler/%s/schedule/%s", scheduler1.SchedulerName, s1.ID())

	tests := []struct {
		schedules         []schedulerSchedules
		expectedSchedules []schedule.Schedule
	}{
		// one version
		{schedulersSchedules(scheduler1), []schedule.Schedule{s1}},
		// two versions same schedule
		{schedulersSchedules(scheduler1bis), []schedule.Schedule{s1, s1}},
		// one version for two schedules
		{schedulersSchedules(scheduler1), []schedule.Schedule{s1}},
	}
	for _, prefix := range []string{"", "/live"} {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
				createSchedulerSchedules(tt.schedules, stores...)

				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, prefix+url, http.NoBody)
				response := executeRequest(router, req)

				checkResponseJSON(t, http.StatusOK, response, toJSON(t, tt.expectedSchedules))
			})
		}
	}
}

// Rule #3: list schedulers response when not found or error: 404 or 500
func TestRestAPIServer_listSchedulers_not_found_or_error(t *testing.T) {
	resolver := slice.NewResolver()

	router := restapi.NewRouter(nil, nil, nil, resolver)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	tests := []struct {
		schedulers       []schedulers.Scheduler
		expectedCode     int
		expectedResponse string
	}{
		{nil, http.StatusOK, "[]"},
		// TODO
		// {[]slice.Scheduler{{SchedulerName: "ERROR"}}, http.StatusInternalServerError, `{"error":"simulated error"}`},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			if tt.schedulers != nil {
				resolver.Reset()
				resolver.Add(tt.schedulers...)
			}
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/schedulers", http.NoBody)
			response := executeRequest(router, req)
			checkResponseJSON(t, tt.expectedCode, response, tt.expectedResponse)
		})
	}
}

// Rule #4: list schedulers response when found: 302 + expected payload
func TestRestAPIServer_listSchedulers_found(t *testing.T) {
	resolver := slice.NewResolver()

	router := restapi.NewRouter(nil, nil, nil, resolver)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	scheduler1 := slice.Scheduler{SchedulerName: "scheduler-1"}
	scheduler2 := slice.Scheduler{SchedulerName: "scheduler-2"}

	tests := []struct {
		schedulers         []slice.Scheduler
		expectedSchedulers []slice.Scheduler
	}{
		{schedulersSlice(scheduler1), schedulersSlice(scheduler1)},
		{schedulersSlice(scheduler1, scheduler2), schedulersSlice(scheduler1, scheduler2)},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			createSchedulers(resolver, tt.expectedSchedulers)

			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/schedulers", http.NoBody)
			response := executeRequest(router, req)
			checkResponseJSON(t, http.StatusOK, response, toJSON(t, tt.schedulers))
		})
	}
}

// Rule #5: search schedules with no search query should return all schedules
func TestRestAPIServer_searchSchedules_no_query(t *testing.T) {
	router, stores, _ := newRouter()

	schedules1 := newSchedules("scheduler-1", 1)
	schedules2 := newSchedules("scheduler-2", 10)

	scheduler1 := schedulerSchedules{
		"scheduler-1",
		schedules1,
	}

	scheduler2 := schedulerSchedules{
		"scheduler-2",
		schedules2,
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	tests := []struct {
		schedules         []schedulerSchedules
		scheduler         string
		query             searchQuery
		expectedCode      int
		expectedSchedules []schedule.Schedule
	}{
		{schedulersSchedules(scheduler1), "scheduler-1", searchQuery{}, http.StatusOK, schedules1},
		{schedulersSchedules(scheduler2), "scheduler-2", searchQuery{}, http.StatusOK, schedules2},
	}

	for _, url := range SchedulesEndpoints {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
				createSchedulerSchedules(tt.schedules, stores...)

				surl := fmt.Sprintf(url, tt.scheduler)
				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, surl, http.NoBody)
				response := executeRequest(router, req)

				checkResponseJSON(t, tt.expectedCode, response, toJSON(t,
					searchResult{
						Found:     len(tt.expectedSchedules),
						Schedules: tt.expectedSchedules,
					}))
			})
		}
	}
}

// Rule #6: search schedules with scheduler name should return corresponding schedules
func TestRestAPIServer_searchSchedules_search_by_schedulerName(t *testing.T) {
	router, stores, _ := newRouter()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	schedules1 := newSchedules("scheduler-1", 1)
	schedules2 := newSchedules("scheduler-2", 10)
	schedules3 := newSchedules("scheduler-3", 10)

	scheduler1 := schedulerSchedules{
		"scheduler-1",
		schedules1,
	}

	scheduler2 := schedulerSchedules{
		"scheduler-2",
		schedules2,
	}

	scheduler3 := schedulerSchedules{
		"scheduler-10",
		schedules3,
	}

	tests := []struct {
		schedules         []schedulerSchedules
		scheduler         string
		expectedCode      int
		expectedSchedules []schedule.Schedule
	}{
		{schedulersSchedules(scheduler1, scheduler2), "scheduler-3", http.StatusOK, []schedule.Schedule{}},
		{schedulersSchedules(scheduler1, scheduler2), "scheduler-1", http.StatusOK, schedules1},
		{schedulersSchedules(scheduler1, scheduler2), "scheduler-2", http.StatusOK, schedules2},
		// exact match: "scheduler-1" query should not return scheduler-10 schedules
		{schedulersSchedules(scheduler1, scheduler2, scheduler3), "scheduler-1", http.StatusOK, schedules1},
	}

	for _, url := range SchedulesEndpoints {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
				createSchedulerSchedules(tt.schedules, stores...)
				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(url, tt.scheduler), http.NoBody)
				response := executeRequest(router, req)

				checkResponseJSON(t, tt.expectedCode, response, toJSON(t, struct {
					Found     int                 `json:"found"`
					Schedules []schedule.Schedule `json:"schedules"`
				}{
					Found:     len(tt.expectedSchedules),
					Schedules: tt.expectedSchedules,
				}))
			})
		}
	}
}

// Rule #7: search schedules with schedule ID should return corresponding schedules
func TestRestAPIServer_searchSchedules_search_by_scheduleID(t *testing.T) {
	router, stores, _ := newRouter()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	schedule1 := newSchedule("scheduler-1", "schedule-1", time.Now())
	scheduler1 := schedulerSchedules{
		"scheduler-1",
		schedulesSlice(schedule1),
	}

	schedule2 := newSchedule("scheduler-2", "schedule-2", time.Now())
	schedule20 := newSchedule("scheduler-2", "schedule-20", time.Now())
	scheduler2 := schedulerSchedules{
		"scheduler-2",
		schedulesSlice(schedule2, schedule20),
	}

	tests := []struct {
		schedules         []schedulerSchedules
		scheduler         string
		query             searchQuery
		expectedCode      int
		expectedSchedules []schedule.Schedule
	}{
		{schedulersSchedules(scheduler1), "scheduler-1", searchQuery{schedulerID: "schedule-0"}, http.StatusOK, []schedule.Schedule{}},
		{schedulersSchedules(scheduler1), "scheduler-1", searchQuery{schedulerID: "schedule-1"}, http.StatusOK, schedulesSlice(schedule1)},
		{schedulersSchedules(scheduler2), "scheduler-2", searchQuery{schedulerID: "schedule-2"}, http.StatusOK, schedulesSlice(schedule2, schedule20)},
		{schedulersSchedules(scheduler2), "scheduler-2", searchQuery{schedulerID: "schedule"}, http.StatusOK, schedulesSlice(schedule2, schedule20)},
	}

	for _, url := range SchedulesEndpoints {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
				createSchedulerSchedules(tt.schedules, stores...)

				surl := fmt.Sprintf(url, tt.scheduler)
				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, tt.query.toURLParams(surl), http.NoBody)
				response := executeRequest(router, req)

				checkResponseJSON(t, tt.expectedCode, response, toJSON(t, struct {
					Found     int                 `json:"found"`
					Schedules []schedule.Schedule `json:"schedules"`
				}{
					Found:     len(tt.expectedSchedules),
					Schedules: tt.expectedSchedules,
				}))
			})
		}
	}
}

// Rule #8: search schedules with max should limit result size, found should match all found in the DB
func TestRestAPIServer_searchSchedules_max(t *testing.T) {
	router, stores, max := newRouter()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	schedules10 := newSchedules("scheduler-1", 10)

	// schedules
	schedulesMax := newSchedules("scheduler-2", max)
	schedulesOverMax := newSchedules("scheduler-3", max+1)

	scheduler1 := schedulerSchedules{
		"scheduler-1",
		schedules10,
	}
	scheduler2 := schedulerSchedules{
		"scheduler-2",
		schedulesMax,
	}
	scheduler3 := schedulerSchedules{
		"scheduler-3",
		schedulesOverMax,
	}

	tests := []struct {
		schedules         []schedulerSchedules
		scheduler         string
		query             searchQuery
		expectedFound     int
		expectedSchedules []schedule.Schedule
	}{
		// max to 1
		{schedulersSchedules(scheduler1), "scheduler-1", searchQuery{max: 1}, 10, schedules10[0:1]},
		// max greater than result set
		{schedulersSchedules(scheduler1), "scheduler-1", searchQuery{max: 20}, 10, schedules10},
		// invalid max, should default to default max
		{schedulersSchedules(scheduler3), "scheduler-3", searchQuery{max: -10}, len(schedulesOverMax), schedulesOverMax[0:max]},
		{schedulersSchedules(scheduler3), "scheduler-3", searchQuery{max: 0}, len(schedulesOverMax), schedulesOverMax[0:max]},
		// no max defined, should default to default max
		{schedulersSchedules(scheduler2), "scheduler-2", searchQuery{}, max, schedulesMax},
		{schedulersSchedules(scheduler3), "scheduler-3", searchQuery{}, len(schedulesOverMax), schedulesOverMax[0:max]},
	}

	for _, url := range SchedulesEndpoints {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
				createSchedulerSchedules(tt.schedules, stores...)

				surl := fmt.Sprintf(url, tt.scheduler)
				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, tt.query.toURLParams(surl), http.NoBody)
				response := executeRequest(router, req)

				checkResponseJSON(t, http.StatusOK, response, toJSON(t, struct {
					Found     int                 `json:"found"`
					Schedules []schedule.Schedule `json:"schedules"`
				}{
					Found:     tt.expectedFound,
					Schedules: tt.expectedSchedules,
				}))
			})
		}
	}
}

// Rule #9: search schedules with sort by, sorting should be applied to the result
func TestRestAPIServer_searchSchedules_sort_by(t *testing.T) {
	router, stores, _ := newRouter()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	now := time.Now()
	s1 := newSchedule("scheduler-1", "schedule-1", now.Add(4*time.Second), now.Add(1*time.Second))
	s2 := newSchedule("scheduler-1", "schedule-2", now.Add(5*time.Second), now.Add(3*time.Second))
	s3 := newSchedule("scheduler-1", "schedule-3", now.Add(2*time.Second), now.Add(2*time.Second))
	s4 := newSchedule("scheduler-1", "schedule-4", now.Add(3*time.Second), now.Add(4*time.Second))
	s5 := newSchedule("scheduler-1", "schedule-5", now.Add(1*time.Second), now.Add(5*time.Second))

	scheduler1 := schedulerSchedules{
		"scheduler-1",
		schedulesSlice(s1, s2, s3, s4, s5),
	}

	schedules := schedulersSchedules(scheduler1)

	tests := []struct {
		schedules         []schedulerSchedules
		query             searchQuery
		expectedSchedules []schedule.Schedule
	}{
		// should default to "timestamp desc"
		{schedules, searchQuery{}, schedulesSlice(s5, s4, s2, s3, s1)},
		// should default to "timestamp desc"
		{schedules, searchQuery{sortOrder: "desc"}, schedulesSlice(s5, s4, s2, s3, s1)},
		// should default to "timestamp asc"
		{schedules, searchQuery{sortOrder: "asc"}, schedulesSlice(s1, s3, s2, s4, s5)},
		// should default to "timestamp desc"
		{schedules, searchQuery{sortField: "timestamp"}, schedulesSlice(s5, s4, s2, s3, s1)},
		{schedules, searchQuery{sortField: "timestamp", sortOrder: "desc"}, schedulesSlice(s5, s4, s2, s3, s1)},
		{schedules, searchQuery{sortField: "timestamp", sortOrder: "asc"}, schedulesSlice(s1, s3, s2, s4, s5)},
		// should default to "id desc"
		{schedules, searchQuery{sortField: "id"}, schedulesSlice(s5, s4, s3, s2, s1)},
		{schedules, searchQuery{sortField: "id", sortOrder: "desc"}, schedulesSlice(s5, s4, s3, s2, s1)},
		{schedules, searchQuery{sortField: "id", sortOrder: "asc"}, schedulesSlice(s1, s2, s3, s4, s5)},
		// should default to "epoch desc"
		{schedules, searchQuery{sortField: "epoch"}, schedulesSlice(s2, s1, s4, s3, s5)},
		{schedules, searchQuery{sortField: "epoch", sortOrder: "desc"}, schedulesSlice(s2, s1, s4, s3, s5)},
		{schedules, searchQuery{sortField: "epoch", sortOrder: "asc"}, schedulesSlice(s5, s3, s4, s1, s2)},
	}
	for _, url := range SchedulesEndpoints {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
				createSchedulerSchedules(tt.schedules, stores...)

				surl := fmt.Sprintf(url, "scheduler-1")
				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, tt.query.toURLParams(surl), http.NoBody)
				response := executeRequest(router, req)

				checkResponseJSON(t, http.StatusOK, response, toJSON(t, struct {
					Found     int                 `json:"found"`
					Schedules []schedule.Schedule `json:"schedules"`
				}{
					Found:     len(tt.expectedSchedules),
					Schedules: tt.expectedSchedules,
				}))
			})
		}
	}
}

// Rule #10: search schedules by epoch range should filter result
func TestRestAPIServer_searchSchedules_search_by_epoch(t *testing.T) {
	router, stores, _ := newRouter()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	now := time.Now()
	s1 := newSchedule("scheduler-1", "schedule-1", now.Add(1*time.Second))
	s2 := newSchedule("scheduler-1", "schedule-2", now.Add(2*time.Second))
	s3 := newSchedule("scheduler-1", "schedule-3", now.Add(3*time.Second))
	s4 := newSchedule("scheduler-1", "schedule-4", now.Add(4*time.Second))
	s5 := newSchedule("scheduler-1", "schedule-5", now.Add(5*time.Second))

	scheduler1 := schedulerSchedules{
		"scheduler-1",
		schedulesSlice(s1, s2, s3, s4, s5),
	}

	schedules := schedulersSchedules(scheduler1)

	tests := []struct {
		schedules         []schedulerSchedules
		query             searchQuery
		expectedCode      int
		expectedSchedules []schedule.Schedule
	}{
		// from only, 0 means none
		{schedules, searchQuery{epochFrom: 0}, http.StatusOK, schedulesSlice(s1, s2, s3, s4, s5)},
		{schedules, searchQuery{epochFrom: 0, epochTo: 0}, http.StatusOK, schedulesSlice(s1, s2, s3, s4, s5)},
		{schedules, searchQuery{epochFrom: 1, epochTo: 0}, http.StatusOK, schedulesSlice(s1, s2, s3, s4, s5)},
		{schedules, searchQuery{epochFrom: 0, epochTo: 1}, http.StatusOK, []schedule.Schedule{}},
		{schedules, searchQuery{epochFrom: now.Add(2 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s2, s3, s4, s5)},
		{schedules, searchQuery{epochFrom: now.Add(5 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s5)},
		{schedules, searchQuery{epochFrom: now.Add(6 * time.Second).Unix()}, http.StatusOK, []schedule.Schedule{}},
		// to only
		{schedules, searchQuery{epochTo: now.Add(1 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s1)},
		{schedules, searchQuery{epochTo: now.Add(3 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s1, s2, s3)},
		{schedules, searchQuery{epochTo: now.Add(5 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s1, s2, s3, s4, s5)},
		{schedules, searchQuery{epochTo: now.Add(6 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s1, s2, s3, s4, s5)},
		{schedules, searchQuery{epochTo: now.Add(-1 * time.Second).Unix()}, http.StatusOK, []schedule.Schedule{}},
		// from and to
		{schedules, searchQuery{epochFrom: 1, epochTo: now.Add(6 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s1, s2, s3, s4, s5)},
		{schedules, searchQuery{epochFrom: now.Add(2 * time.Second).Unix(), epochTo: now.Add(4 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s2, s3, s4)},
		{schedules, searchQuery{epochFrom: now.Add(1 * time.Second).Unix(), epochTo: now.Add(6 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s1, s2, s3, s4, s5)},
		{schedules, searchQuery{epochFrom: now.Add(-10 * time.Second).Unix(), epochTo: now.Add(10 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s1, s2, s3, s4, s5)},
		// if to > from, only from is taken into account
		{schedules, searchQuery{epochFrom: now.Add(3 * time.Second).Unix(), epochTo: now.Add(2 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s3, s4, s5)},
		{schedules, searchQuery{epochFrom: now.Add(0 * time.Second).Unix(), epochTo: now.Add(-2 * time.Second).Unix()}, http.StatusOK, schedulesSlice(s1, s2, s3, s4, s5)},
	}
	for _, url := range SchedulesEndpoints {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
				createSchedulerSchedules(tt.schedules, stores...)

				surl := fmt.Sprintf(url, "scheduler-1")
				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, tt.query.toURLParams(surl), http.NoBody)
				response := executeRequest(router, req)

				checkResponseJSON(t, tt.expectedCode, response, toJSON(t, struct {
					Found     int                 `json:"found"`
					Schedules []schedule.Schedule `json:"schedules"`
				}{
					Found:     len(tt.expectedSchedules),
					Schedules: tt.expectedSchedules,
				}))
			})
		}
	}
}

// Rule #11: search schedules by multiple criteria
func TestRestAPIServer_live_searchSchedules_search_multicriteria(t *testing.T) {
	router, stores, _ := newRouter()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	now := time.Now()
	s1 := newSchedule("scheduler-1", "schedule-1", now.Add(1*time.Second))
	s2 := newSchedule("scheduler-1", "schedule-2", now.Add(2*time.Second))
	s3 := newSchedule("scheduler-2", "schedule-3", now.Add(3*time.Second))
	s4 := newSchedule("scheduler-2", "schedule-4", now.Add(4*time.Second))
	s5 := newSchedule("scheduler-2", "schedule-5", now.Add(5*time.Second))

	scheduler1 := schedulerSchedules{
		"scheduler-1",
		schedulesSlice(s1, s2),
	}
	scheduler2 := schedulerSchedules{
		"scheduler-2",
		schedulesSlice(s3, s4, s5),
	}

	schedules := schedulersSchedules(scheduler1, scheduler2)

	tests := []struct {
		schedules         []schedulerSchedules
		scheduler         string
		query             searchQuery
		expectedCode      int
		expectedFound     int
		expectedSchedules []schedule.Schedule
	}{
		{schedules, "scheduler-1", searchQuery{schedulerName: scheduler1.SchedulerName, schedulerID: s1.ID(), max: 1}, http.StatusOK, 1, schedulesSlice(s1)},
		{schedules, "scheduler-2", searchQuery{schedulerName: scheduler2.SchedulerName, schedulerID: "sch", max: 2, sortField: "id", sortOrder: "desc"}, http.StatusOK, 3, schedulesSlice(s5, s4)},
		{schedules, "scheduler-2", searchQuery{schedulerName: scheduler2.SchedulerName, schedulerID: "sch", max: 10, sortField: "id", sortOrder: "desc"}, http.StatusOK, 3, schedulesSlice(s5, s4, s3)},
		{schedules, "scheduler-2", searchQuery{schedulerName: scheduler2.SchedulerName, schedulerID: "sch", max: -1, sortField: "id", sortOrder: "desc"}, http.StatusOK, 3, schedulesSlice(s5, s4, s3)},
		{schedules, "scheduler-2", searchQuery{schedulerName: scheduler2.SchedulerName, max: -1, sortField: "id", sortOrder: "desc"}, http.StatusOK, 3, schedulesSlice(s5, s4, s3)},
		{schedules, "scheduler-2", searchQuery{schedulerName: scheduler2.SchedulerName, max: -1, sortField: "id", sortOrder: "asc"}, http.StatusOK, 3, schedulesSlice(s3, s4, s5)},
		{schedules, "scheduler-2", searchQuery{schedulerName: scheduler2.SchedulerName, max: 2, epochFrom: s3.Epoch(), epochTo: s4.Epoch(), sortField: "id", sortOrder: "asc"}, http.StatusOK, 2, schedulesSlice(s3, s4)},
		{schedules, "scheduler-2", searchQuery{schedulerName: scheduler2.SchedulerName, max: 2, epochFrom: s3.Epoch(), epochTo: s4.Epoch(), sortField: "id", sortOrder: "desc"}, http.StatusOK, 2, schedulesSlice(s4, s3)},
		{schedules, "scheduler-2", searchQuery{schedulerName: scheduler2.SchedulerName, max: 1, epochFrom: s3.Epoch(), epochTo: s4.Epoch(), sortField: "id", sortOrder: "asc"}, http.StatusOK, 2, schedulesSlice(s3)},
		{schedules, "scheduler-2", searchQuery{schedulerName: scheduler2.SchedulerName, max: 99, epochFrom: s3.Epoch(), epochTo: s4.Epoch(), sortField: "id", sortOrder: "desc"}, http.StatusOK, 2, schedulesSlice(s4, s3)},
	}

	for _, url := range SchedulesEndpoints {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("case #%v (%s)", i+1, url), func(t *testing.T) {
				createSchedulerSchedules(tt.schedules, stores...)

				surl := fmt.Sprintf(url, tt.scheduler)
				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, tt.query.toURLParams(surl), http.NoBody)
				response := executeRequest(router, req)

				checkResponseJSON(t, tt.expectedCode, response, toJSON(t, struct {
					Found     int                 `json:"found"`
					Schedules []schedule.Schedule `json:"schedules"`
				}{
					Found:     tt.expectedFound,
					Schedules: tt.expectedSchedules,
				}))
			})
		}
	}
}

// Rule #12: stats endpoint should return a list of scheduler with metrics
func TestRestAPIServer_stats(t *testing.T) {
	resolver := slice.NewResolver()
	resolver.Add(slice.Scheduler{SchedulerName: "scheduler-1"})
	resolver.Add(slice.Scheduler{SchedulerName: "scheduler-2"})

	router, stores, _ := newRouter(resolver)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	schedule1 := newSchedule("scheduler-1", "schedule-1", time.Now())
	scheduler1 := schedulerSchedules{
		"scheduler-1",
		schedulesSlice(schedule1),
	}

	schedule2 := newSchedule("scheduler-2", "schedule-2", time.Now())
	schedule3 := newSchedule("scheduler-2", "schedule-3", time.Now())
	scheduler2 := schedulerSchedules{
		"scheduler-2",
		schedulesSlice(schedule2, schedule3),
	}

	tests := []struct {
		schedules []schedulerSchedules
	}{
		{schedulersSchedules(scheduler1, scheduler2)},
	}

	for _, url := range []string{"/stats"} {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
				createSchedulerSchedules(tt.schedules, stores...)

				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
				response := executeRequest(router, req)

				checkResponseJSON(t, http.StatusOK, response, toJSON(t, []struct {
					SchedulerName string `json:"scheduler"`
					TotalLive     int    `json:"total_live"`
					TotalHistory  int    `json:"total_history"`
					Total         int    `json:"total"`
				}{{
					"scheduler-1",
					1,
					1,
					1,
				}, {
					"scheduler-2",
					2,
					2,
					2,
				}}))
			})
		}
	}
}
