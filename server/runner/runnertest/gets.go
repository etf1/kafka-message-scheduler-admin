package runnertest

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	log "github.com/sirupsen/logrus"
)

func getSchedulers(timeout time.Duration) (resp *http.Response, err error) {
	return get("/schedulers", timeout)
}

func getSchedules(schedulerName string, timeout time.Duration) (resp *http.Response, err error) {
	return get(fmt.Sprintf("/scheduler/%s/schedules", schedulerName), timeout)
}
func getSchedule(schedulerName, id string, timeout time.Duration) (resp *http.Response, err error) {
	return get(fmt.Sprintf("/scheduler/%s/schedule/%s", schedulerName, id), timeout)
}

func getLiveSchedules(schedulerName string, timeout time.Duration) (resp *http.Response, err error) {
	return get(fmt.Sprintf("/live/scheduler/%s/schedules", schedulerName), timeout)
}
func getLiveSchedule(schedulerName, id string, timeout time.Duration) (resp *http.Response, err error) {
	return get(fmt.Sprintf("/live/scheduler/%s/schedule/%s", schedulerName, id), timeout)
}

func get(path string, timeout time.Duration) (*http.Response, error) {
	addr := config.ServerAddr()

	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}

	if !config.APIServerOnly() {
		path = "/api" + path
	}

	url := "http://" + addr + path

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: timeout,
	}

	log.Printf("calling get url (2): %v", url)

	return client.Do(req)
}
