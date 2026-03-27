package mini_test

import (
	"log"
	"testing"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner/mini"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner/runnertest"
)

// runMiniTest starts a mini runner on a free port, registers cleanup via t.Cleanup,
// waits for the server to be ready, and calls checkFn to exercise the endpoint.
func runMiniTest(t *testing.T, checkFn func() error) {
	t.Helper()

	runner := mini.NewRunner()
	config.SetServerAddr(helper.NextServerAddr("localhost"))

	exitchan := make(chan bool, 1)
	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to start the mini runner: %v", err)
		}
		exitchan <- true
	}()

	t.Cleanup(func() {
		runner.Close()
		<-exitchan
	})

	if err := helper.WaitForHTTPServer(config.ServerAddr()); err != nil {
		t.Fatalf("unreachable host: %v", err)
	}

	if err := checkFn(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// Rule #1: runner must expose the api server endpoint /schedulers
func TestMiniRunner_schedulers(t *testing.T) {
	runMiniTest(t, runnertest.CheckSchedulersEndPoint)
}

// Rule #2: runner must expose the api server endpoint /schedules
func TestMiniRunner_schedules(t *testing.T) {
	runMiniTest(t, func() error {
		return runnertest.CheckSchedulesEndPoint("scheduler-1")
	})
}

// Rule #3: runner must expose the api server endpoint /scheduler/{name}/schedule/{id}
func TestMiniRunner_schedule_detail(t *testing.T) {
	runMiniTest(t, func() error {
		return runnertest.CheckScheduleDetailEndPoint("scheduler-1", "schedule-1")
	})
}

// Rule #4: runner must expose the api server endpoint /live/schedules
func TestMiniRunner_live_schedules(t *testing.T) {
	runMiniTest(t, func() error {
		return runnertest.CheckLiveSchedulesEndPoint("scheduler-1")
	})
}

// Rule #5: runner must expose the api server endpoint /live/scheduler/{name}/schedule/{id}
func TestMiniRunner_detail_live_schedule(t *testing.T) {
	runMiniTest(t, func() error {
		return runnertest.CheckLiveScheduleDetailEndPoint("scheduler-1", "schedule-1")
	})
}
