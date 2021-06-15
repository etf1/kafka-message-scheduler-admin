// INTEGRATION TESTS
package kafka_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner/kafka"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner/runnertest"
)

func TestKafkaRunner_schedulers(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	exitchan := make(chan bool, 1)

	dataDir := "./" + helper.GenRandString("db-")
	defer func() {
		os.RemoveAll(dataDir)
	}()

	runner := kafka.NewRunner(dataDir)

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

func TestKafkaRunner_schedules(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	exitchan := make(chan bool, 1)

	dataDir := "./" + helper.GenRandString("db-")
	defer func() {
		os.RemoveAll(dataDir)
	}()

	runner := kafka.NewRunner(dataDir)

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := runnertest.CheckSchedulesEndPoint("scheduler", runner, exitchan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKafkaRunner_schedule_detail(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	exitchan := make(chan bool, 1)

	dataDir := "./" + helper.GenRandString("db-")
	defer func() {
		os.RemoveAll(dataDir)
	}()

	runner := kafka.NewRunner(dataDir)

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := runnertest.CheckScheduleDetailEndPoint("localhost", "schedule-1", runner, exitchan)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestKafkaRunner_live_schedules(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	dataDir := "./" + helper.GenRandString("db-")
	defer func() {
		os.RemoveAll(dataDir)
	}()

	exitchan := make(chan bool, 1)

	runner := kafka.NewRunner(dataDir)

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := runnertest.CheckLiveSchedulesEndPoint("localhost", runner, exitchan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKafkaRunner_live_schedule_detail(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	exitchan := make(chan bool, 1)

	dataDir := "./" + helper.GenRandString("db-")
	defer func() {
		os.RemoveAll(dataDir)
	}()

	runner := kafka.NewRunner(dataDir)

	go func() {
		if err := runner.Start(); err != nil {
			log.Printf("failed to create the default kafka runner: %v", err)
		}
		exitchan <- true
	}()

	// wait for the goroutine to be scheduled
	time.Sleep(1 * time.Second)

	err := runnertest.CheckLiveScheduleDetailEndPoint("localhost", "schedule-1", runner, exitchan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
