package mini_test

import (
	"log"
	"testing"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/runner/mini"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner/runnertest"
)

// Rule #1: mini runner must expose the api server endpoint /schedulers
func TestMiniRunner_schedulers(t *testing.T) {
	exitchan := make(chan bool, 1)

	runner := mini.NewRunner()

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := runnertest.CheckSchedulersEndPoint(runner, exitchan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Rule #2: mini runner must expose the api server endpoint /schedules
func TestMiniRunner_schedules(t *testing.T) {
	exitchan := make(chan bool, 1)

	runner := mini.NewRunner()

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := runnertest.CheckSchedulesEndPoint("scheduler-1", runner, exitchan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Rule #3: mini runner must expose the api server endpoint /scheduler/{name}/schedule/{id}
func TestMiniRunner_schedule_detail(t *testing.T) {
	exitchan := make(chan bool, 1)

	runner := mini.NewRunner()

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := runnertest.CheckScheduleDetailEndPoint("scheduler-1", "schedule-1", runner, exitchan)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// Rule #4: mini runner must expose the api server endpoint /live/schedules
func TestMiniRunner_live_schedules(t *testing.T) {
	exitchan := make(chan bool, 1)

	runner := mini.NewRunner()

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := runnertest.CheckLiveSchedulesEndPoint("scheduler-1", runner, exitchan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Rule #5: mini runner must expose the api server endpoint /live/scheduler/{name}/schedule/{id}
func TestMiniRunner_detail_live_schedule(t *testing.T) {
	exitchan := make(chan bool, 1)

	runner := mini.NewRunner()

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := runnertest.CheckLiveScheduleDetailEndPoint("scheduler-1", "schedule-1", runner, exitchan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
