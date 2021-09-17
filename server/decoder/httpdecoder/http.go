package httpdecoder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/store/rest"
	"github.com/etf1/kafka-message-scheduler/schedule"
	"github.com/etf1/kafka-message-scheduler/schedule/kafka"
	log "github.com/sirupsen/logrus"
)

var (
	ErrUnknownScheduleType = fmt.Errorf("unknown schedule type")
)

type Decoder struct {
	URL string
}

// Post sends an http post request with the payload and return the response
func (h Decoder) Post(payload interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, h.URL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Println("response status code:", resp.StatusCode)
	log.Println("response headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %v %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Decode returns a copy of the input schedule, its field 'value' replaced by the response from the http request
func (h Decoder) Decode(s schedule.Schedule) (schedule.Schedule, error) {
	switch sch := s.(type) {
	case *kafka.Schedule:
		if len(sch.Message.Value) == 0 {
			return s, nil
		}

		data, err := h.Post(sch)
		if err != nil {
			return s, err
		}

		if len(data) != 0 {
			sch.Message.Value = data
		}

		return sch, nil
	case rest.Schedule:
		if len(sch.MessageValue) == 0 {
			return s, nil
		}

		data, err := h.Post(sch)
		if err != nil {
			return s, err
		}

		if len(data) != 0 {
			sch.MessageValue = data
		}

		return sch, nil
	default:
		return s, fmt.Errorf("%w: %T", ErrUnknownScheduleType, s)
	}
}
