package runnertest

import (
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	timeout = 5 * time.Second
)

type Closable interface {
	Close()
}

func CheckResponse(resp *http.Response, err error) error {
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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

func CheckSchedulersEndPoint() error {
	return CheckResponse(getSchedulers(timeout))
}

func CheckSchedulesEndPoint(schedulerName string) error {
	return CheckResponse(getSchedules(schedulerName, timeout))
}

func CheckScheduleDetailEndPoint(schedulerName, scheduleID string) error {
	return CheckResponse(getSchedule(schedulerName, scheduleID, timeout))
}

func CheckLiveSchedulesEndPoint(schedulerName string) error {
	return CheckResponse(getLiveSchedules(schedulerName, timeout))
}

func CheckLiveScheduleDetailEndPoint(schedulerName, scheduleID string) error {
	return CheckResponse(getLiveSchedule(schedulerName, scheduleID, timeout))
}
