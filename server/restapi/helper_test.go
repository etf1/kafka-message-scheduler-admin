package restapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	stdsort "sort"
	"testing"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/db/simple"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/slice"
	"github.com/etf1/kafka-message-scheduler-admin/server/restapi"
	"github.com/etf1/kafka-message-scheduler-admin/server/sort"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/hmap"

	"github.com/etf1/kafka-message-scheduler/schedule"
	simple_schedule "github.com/etf1/kafka-message-scheduler/schedule/simple"
)

const (
	SchedulesEndpoint     = "/scheduler/%s/schedules"
	LiveSchedulesEndpoint = "/live/scheduler/%s/schedules"
)

var (
	simpleSchedule    = simple_schedule.NewSchedule
	emptySearchResult = searchResult{
		Found:     0,
		Schedules: []schedule.Schedule{},
	}
)

type searchResult struct {
	Found     int                 `json:"found"`
	Schedules []schedule.Schedule `json:"schedules"`
}

type schedulerSchedules struct {
	SchedulerName string              `json:"scheduler"`
	Schedules     []schedule.Schedule `json:"schedules"`
}

func sortSchedules(arr []schedule.Schedule, sb ...sort.SortBy) []schedule.Schedule {
	result := append([]schedule.Schedule{}, arr...)
	sortBy := sort.DefaultSortBy
	if len(sb) > 0 {
		sortBy = sb[0]
	}
	stdsort.Sort(sort.NewSort(result, sortBy))
	return result
}

func schedulesSlice(arr ...schedule.Schedule) []schedule.Schedule {
	return arr
}

func schedulersSlice(arr ...slice.Scheduler) []slice.Scheduler {
	return arr
}

func schedulersSchedules(arr ...schedulerSchedules) []schedulerSchedules {
	return arr
}

type searchQuery struct {
	schedulerName string
	schedulerID   string
	max           int
	epochFrom     int64
	epochTo       int64
	sortField     string
	sortOrder     string
}

func (s searchQuery) toURLParams(base string) string {
	v := url.Values{}
	if s.schedulerName != "" {
		v.Set("scheduler-name", s.schedulerName)
	}
	if s.schedulerID != "" {
		v.Set("schedule-id", s.schedulerID)
	}
	if s.epochFrom != 0 {
		v.Set("epoch-from", fmt.Sprint(s.epochFrom))
	}
	if s.epochTo != 0 {
		v.Set("epoch-to", fmt.Sprint(s.epochTo))
	}
	if s.sortField != "" || s.sortOrder != "" {
		sortBy := ""
		if s.sortField != "" {
			sortBy = s.sortField
		}
		if s.sortOrder != "" {
			if s.sortField != "" {
				sortBy += " "
			}
			sortBy += s.sortOrder
		}
		v.Set("sort-by", sortBy)
	}
	if s.max != 0 {
		v.Set("max", fmt.Sprint(s.max))
	}
	res := base
	if encoded := v.Encode(); encoded != "" {
		res += "?" + encoded
	}

	return res
}

func createSchedulerSchedules(schedulers []schedulerSchedules, stores ...*hmap.Hmap) {
	for _, st := range stores {
		st.Reset()
		for _, scheduler := range schedulers {
			for _, sch := range scheduler.Schedules {
				switch v := sch.(type) {
				case store.Schedule:
					st.Add(v.SchedulerName, v.Schedule)
				default:
					st.Add(scheduler.SchedulerName, v)
				}
			}
		}
	}
}

func createSchedules(schedulerName string, schedules []schedule.Schedule, m *hmap.Hmap) {
	m.Reset()
	m.Add(schedulerName, schedules...)
}

func createSchedulers(service *slice.Slice, schs []slice.Scheduler) {
	service.Reset()

	arr := make([]schedulers.Scheduler, len(schs))
	for i, s := range schs {
		arr[i] = s
	}

	service.Add(arr...)
}

var index int = 0

func newSchedules(schedulerName string, size int, sb ...sort.SortBy) []schedule.Schedule {
	result := make([]schedule.Schedule, size)
	for i := 0; i < size; i++ {
		index++
		t := time.Now().Add(time.Duration(index) * time.Second)
		result[i] = store.Schedule{
			SchedulerName: schedulerName,
			Schedule:      simpleSchedule(fmt.Sprintf("schedule-%v", index), t.Unix(), t),
		}
	}

	return sortSchedules(result, sb...)
}

func newSchedule(schedulerName, scheduleID string, epoch interface{}, timestamp ...time.Time) schedule.Schedule {
	return store.Schedule{
		SchedulerName: schedulerName,
		Schedule:      simpleSchedule(scheduleID, epoch, timestamp...),
	}
}

func newServer() (restapi.Server, []*hmap.Hmap) {
	cold := hmap.NewStore()
	live := hmap.NewStore()
	srv := restapi.NewServer(simple.DB{
		Store: cold,
	}, simple.DB{
		Store: live,
	}, nil)

	return srv, []*hmap.Hmap{cold, live}
}

func executeRequest(router http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response code %d, got %d\n", expected, actual)
	}
}

func ifThenElse(cond bool, a interface{}, b interface{}) interface{} {
	if cond {
		return a
	}
	return b
}

func toJson(t *testing.T, v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(b) == "null" {
		return nil
	}
	return b
}

func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	if s1 == "" && s2 != "" {
		return false, nil
	}
	if s1 != "" && s2 == "" {
		return false, nil
	}
	if s1 == "" && s2 == "" {
		return true, nil
	}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("error marshalling s1 : %s", err)
	}

	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("error marshalling s2: %s", err)
	}

	return reflect.DeepEqual(o1, o2), nil
}

func checkResponseJson(t *testing.T, expectedCode int, response *httptest.ResponseRecorder, expected interface{}) {
	checkResponseCode(t, expectedCode, response.Code)

	var expectedResponse string

	switch v := expected.(type) {
	case string:
		expectedResponse = v
	case []byte:
		expectedResponse = string(v)
	}

	bodyString := response.Body.String()

	t.Logf("body=%v", bodyString)

	// if bodyString != expectedResponse {
	// 	t.Errorf("unexpected body: %v, expected: %s", bodyString, expected)
	// }

	eq, err := AreEqualJSON(bodyString, expectedResponse)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !eq {
		t.Errorf("unexpected body: %v, expected: %s", bodyString, expected)
	}
}
