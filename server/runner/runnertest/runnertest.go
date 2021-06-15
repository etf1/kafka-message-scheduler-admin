package runnertest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Closable interface {
	Close()
}

func CheckSchedulersEndPoint(runner Closable, exitchan chan bool) error {
	defer runner.Close()

	resp, err := getSchedulers(1 * time.Second)
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	log.Printf("http response body: %s", body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	log.Printf("http status code: %v", resp.Status)

loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			runner.Close()
		case <-exitchan:
			break loop
		}
	}

	return nil
}

func CheckSchedulesEndPoint(schedulerName string, runner Closable, exitchan chan bool) error {
	defer runner.Close()

	resp, err := getSchedules(schedulerName, 1*time.Second)
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	log.Printf("http response body: %s", body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	log.Printf("http status code: %v", resp.Status)

loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			runner.Close()
		case <-exitchan:
			break loop
		}
	}

	return nil
}

func CheckScheduleDetailEndPoint(schedulerName, scheduleID string, runner Closable, exitchan chan bool) error {
	defer runner.Close()

	resp, err := getSchedule(schedulerName, scheduleID, 1*time.Second)
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	log.Printf("http response body: %s", body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	log.Printf("http status code: %v", resp.Status)

loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			runner.Close()
		case <-exitchan:
			break loop
		}
	}
	return nil
}

func CheckLiveSchedulesEndPoint(schedulerName string, runner Closable, exitchan chan bool) error {
	defer runner.Close()

	resp, err := getLiveSchedules(schedulerName, 1*time.Second)
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	log.Printf("http response body: %s", body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	log.Printf("http status code: %v", resp.Status)

loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			runner.Close()
		case <-exitchan:
			break loop
		}
	}
	return nil
}

func CheckLiveScheduleDetailEndPoint(schedulerName, scheduleID string, runner Closable, exitchan chan bool) error {
	defer runner.Close()

	resp, err := getLiveSchedule(schedulerName, scheduleID, 1*time.Second)
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	log.Printf("http response body: %s", body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	log.Printf("http status code: %v", resp.Status)

loop:
	for {
		select {
		case <-time.After(2 * time.Second):
			runner.Close()
		case <-exitchan:
			break loop
		}
	}
	return nil
}
