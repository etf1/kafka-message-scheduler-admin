// INTEGRATION TESTS
package kafka_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner/kafka"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner/runnertest"
)

// Rule #1: runner must expose the api server endpoint /schedulers
func TestKafkaRunner_schedulers(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	exitchan := make(chan bool, 1)

	dataDir := "./" + helper.GenRandString("db-")
	defer func() {
		os.RemoveAll(dataDir)
	}()

	runner := kafka.NewRunner(dataDir)

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
			log.Printf("test: closing runner")
			runner.Close()
		case <-exitchan:
			log.Printf("test: received exit chan")
			break loop
		}
	}
}

// Rule #2: runner must expose the api server endpoint /schedules
func TestKafkaRunner_schedules(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	exitchan := make(chan bool, 1)

	dataDir := "./" + helper.GenRandString("db-")
	defer func() {
		os.RemoveAll(dataDir)
	}()

	runner := kafka.NewRunner(dataDir)

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

	err = runnertest.CheckSchedulesEndPoint("scheduler")
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
func TestKafkaRunner_schedule_detail(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	exitchan := make(chan bool, 1)

	dataDir := "./" + helper.GenRandString("db-")
	defer func() {
		os.RemoveAll(dataDir)
	}()

	runner := kafka.NewRunner(dataDir)

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

	err = runnertest.CheckScheduleDetailEndPoint("localhost", "schedule-1")
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
func TestKafkaRunner_live_schedules(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	dataDir := "./" + helper.GenRandString("db-")
	defer func() {
		os.RemoveAll(dataDir)
	}()

	exitchan := make(chan bool, 1)

	runner := kafka.NewRunner(dataDir)

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

	err = runnertest.CheckLiveSchedulesEndPoint("localhost")
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
func TestKafkaRunner_live_schedule_detail(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	exitchan := make(chan bool, 1)

	dataDir := "./" + helper.GenRandString("db-")
	defer func() {
		os.RemoveAll(dataDir)
	}()

	runner := kafka.NewRunner(dataDir)

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

	err = runnertest.CheckLiveScheduleDetailEndPoint("localhost", "schedule-1")
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
