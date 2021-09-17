// INTEGRATION TESTS

// requirements: a running kafka-message-scheduler:mini
// you can start the integration stack with: `make up`
// or `docker-compose -p dev restart scheduler`
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

// Rule #1: Get should work as expected
func TestRest_Get(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	schedulerName := getSchedulerName()

	dec := &helper.KafkaMessageSimpleDecoder{}
	rstore := rest.NewStore(
		httpresolver.Resolver{
			Hosts: []string{schedulerName},
		},
		dec,
	)

	result, err := rstore.Get(schedulerName, "schedule-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) == 0 {
		t.Errorf("unexpected result: %v", result)
	}

	if dec.Called == 0 {
		t.Errorf("unexpected decoder count: %v", dec.Called)
	}
}

// Rule #2: List should work as expected
func TestRest_List(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	schedulerName := getSchedulerName()

	dec := &helper.KafkaMessageSimpleDecoder{}
	rstore := rest.NewStore(
		httpresolver.Resolver{
			Hosts: []string{schedulerName},
		},
		dec,
	)

	// wait for goroutines to be scheduled
	time.Sleep(1 * time.Second)

	result, err := rstore.List(schedulerName)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
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
		t.Errorf("unexpected schedules length: %v", v)
	}
	if dec.Called == 0 {
		t.Errorf("unexpected decoder count: %v", dec.Called)
	}
}
