package mini_test

import (
	"log"
	"testing"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner/mini"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner/runnertest"
)

// Rule #1: runner must expose the api server endpoint /schedulers
func TestMiniRunner_schedulers(t *testing.T) {
	exitchan := make(chan bool, 1)

	runner := mini.NewRunner()

	// set a random port to avoid conflict
	config.SetServerAddr(helper.NextServerAddr("localhost"))

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := helper.WaitForHTTPServer(config.ServerAddr())
	if err != nil {
		t.Errorf("unreachable host: %v", err)
	}

	err = runnertest.CheckSchedulersEndPoint()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// wait for previous goroutine to exit
loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			runner.Close()
		case <-exitchan:
			break loop
		}
	}
}

// Rule #2: runner must expose the api server endpoint /schedules
func TestMiniRunner_schedules(t *testing.T) {
	exitchan := make(chan bool, 1)

	runner := mini.NewRunner()

	// set a random port to avoid conflict
	config.SetServerAddr(helper.NextServerAddr("localhost"))

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := helper.WaitForHTTPServer(config.ServerAddr())
	if err != nil {
		t.Errorf("unreachable host: %v", err)
	}

	err = runnertest.CheckSchedulesEndPoint("scheduler-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// wait for previous goroutine to exit
loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			runner.Close()
		case <-exitchan:
			break loop
		}
	}
}

// Rule #3: runner must expose the api server endpoint /scheduler/{name}/schedule/{id}
func TestMiniRunner_schedule_detail(t *testing.T) {
	exitchan := make(chan bool, 1)

	runner := mini.NewRunner()

	// set a random port to avoid conflict
	config.SetServerAddr(helper.NextServerAddr("localhost"))

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := helper.WaitForHTTPServer(config.ServerAddr())
	if err != nil {
		t.Errorf("unreachable host: %v", err)
	}

	err = runnertest.CheckScheduleDetailEndPoint("scheduler-1", "schedule-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// wait for previous goroutine to exit
loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			runner.Close()
		case <-exitchan:
			break loop
		}
	}
}

// Rule #4: runner must expose the api server endpoint /live/schedules
func TestMiniRunner_live_schedules(t *testing.T) {
	exitchan := make(chan bool, 1)

	runner := mini.NewRunner()

	// set a random port to avoid conflict
	config.SetServerAddr(helper.NextServerAddr("localhost"))

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := helper.WaitForHTTPServer(config.ServerAddr())
	if err != nil {
		t.Errorf("unreachable host: %v", err)
	}

	err = runnertest.CheckLiveSchedulesEndPoint("scheduler-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// wait for previous goroutine to exit
loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			runner.Close()
		case <-exitchan:
			break loop
		}
	}
}

// Rule #5: runner must expose the api server endpoint /live/scheduler/{name}/schedule/{id}
func TestMiniRunner_detail_live_schedule(t *testing.T) {
	exitchan := make(chan bool, 1)

	runner := mini.NewRunner()

	// set a random port to avoid conflict
	config.SetServerAddr(helper.NextServerAddr("localhost"))

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := helper.WaitForHTTPServer(config.ServerAddr())
	if err != nil {
		t.Errorf("unreachable host: %v", err)
	}

	err = runnertest.CheckLiveScheduleDetailEndPoint("scheduler-1", "schedule-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// wait for previous goroutine to exit
loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			runner.Close()
		case <-exitchan:
			break loop
		}
	}
}
