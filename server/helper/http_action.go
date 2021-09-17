package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	DefaultTimeout = 5 * time.Second
	PortRange      = 500
	PortStartRange = 9002
	MaxRetries     = 5
)

var (
	mu sync.Mutex
)

func Get(host, url string, timeout time.Duration) (*http.Response, error) {
	full := "http://" + host + url
	log.Printf("calling get url: %v", full)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, full, nil)

	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: timeout,
	}

	return client.Do(req)
}

type CheckResponse func(*http.Response) error

var CheckResponseStatusOK = CheckResponse(func(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %v", resp.StatusCode)
	}
	return nil
})

var CheckResponseNil = CheckResponse(func(resp *http.Response) error {
	return nil
})

func DecodeJSON(host, url string, timeout time.Duration, v interface{}, checkResponse ...CheckResponse) error {
	resp, err := Get(host, url, DefaultTimeout)
	if err != nil {
		return fmt.Errorf("cannot get info from host %v: %v", host, err)
	}
	defer resp.Body.Close()

	check := CheckResponseStatusOK
	if len(checkResponse) > 0 {
		check = checkResponse[0]
	}

	err = check(resp)
	if err != nil {
		return fmt.Errorf("invalid response from host %v: %w", host, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot readall of the body from host %v: %v", host, resp.StatusCode)
	}
	log.Printf("http response body: %s", body)

	err = json.NewDecoder(bytes.NewReader(body)).Decode(v)
	if err != nil {
		return fmt.Errorf("unable to unmarshall json body: %v", err)
	}

	return nil
}

func WaitForHTTPServer(addr string) error {
	count := 1

	for {
		timeout := 1 * time.Second

		_, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			log.Printf("unreachable host %v: %v", addr, err)
		} else {
			log.Printf("reachable host %v", addr)
			return nil
		}

		time.Sleep(timeout)

		count++
		if count == MaxRetries {
			return fmt.Errorf("unreachable host after %v retries: %v", MaxRetries, addr)
		}
	}
}

func NextServerAddr(prefix string) string {
	defer mu.Unlock()
	mu.Lock()
	// TODO: check is the port is available before returning
	return fmt.Sprintf("%s:%d", prefix, PortStartRange+RandNumWithMax(PortRange))
}
