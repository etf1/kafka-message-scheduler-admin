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

func CheckSchedulersEndPoint() error {
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	log.Printf("http status code: %v", resp.Status)

	return nil
}

func CheckSchedulesEndPoint(schedulerName string) error {
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

	return nil
}

func CheckScheduleDetailEndPoint(schedulerName, scheduleID string) error {
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

	return nil
}

func CheckLiveSchedulesEndPoint(schedulerName string) error {
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

	return nil
}

func CheckLiveScheduleDetailEndPoint(schedulerName, scheduleID string) error {
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

	return nil
}
