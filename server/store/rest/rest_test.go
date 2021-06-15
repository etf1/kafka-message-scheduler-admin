// INTEGRATION TESTS

// requirements: a running kafka-message-scheduler:mini
// you can start the integration stack with: `make dev.up`
// WARNING: schedules defined in scheduler made be gone when triggered
// you should then restart the scheduler service: `make dev.restart.scheduler`

package rest_test

import (
	"testing"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/rest"
)

func getSchedulerName() string {
	if helper.IsRunningInDocker() {
		return "scheduler"
	}
	return "localhost"
}

func TestRest_Get(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	schedulerName := getSchedulerName()

	rstore := rest.NewStore(
		httpresolver.Resolver{
			Hosts: []string{schedulerName},
		},
	)

	result, err := rstore.Get(schedulerName, "schedule-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) == 0 {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestRest_List(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	schedulerName := getSchedulerName()

	rstore := rest.NewStore(
		httpresolver.Resolver{
			Hosts: []string{schedulerName},
		},
	)

	// wait for goroutines to be scheduled
	time.Sleep(1 * time.Second)

	result, err := rstore.List(schedulerName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count := 0
	for {
		_, ok := <-result
		if !ok {
			break
		}
		count++
	}

	if v := count; v == 0 {
		t.Fatalf("unexpected schedules length: %v", v)
	}
}
