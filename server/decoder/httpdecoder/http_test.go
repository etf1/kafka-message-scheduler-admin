// INTEGRATION TESTS
package httpdecoder_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/etf1/kafka-message-scheduler-admin/server/decoder/httpdecoder"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler/schedule/kafka"
)

// Rule #1: when received http response code != 200, the kafka message body should be unchanged
func TestHTTPDecoder_not200(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)
	tests := []struct {
		code int
	}{
		{http.StatusNotFound},
		{http.StatusBadRequest},
		{http.StatusInternalServerError},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("http_#%v", i), func(t *testing.T) {
			func() {
				server := helper.MockServer(tt.code, "response")
				defer server.Close()

				dec := httpdecoder.Decoder{
					URL: server.URL,
				}

				topic := "topic"
				sch := kafka.Schedule{
					Message: &confluent.Message{
						TopicPartition: confluent.TopicPartition{
							Topic: &topic,
						},
						Key:   []byte("video-1"),
						Value: []byte("content"),
					},
				}

				// create a copy of the schedule
				sch2 := helper.CopyKafkaSchedule(sch)

				sch3, err := dec.Decode(&sch)
				if err == nil {
					t.Error("expected error")
				}

				sch4, ok := sch3.(*kafka.Schedule)
				if !ok {
					t.Errorf("unexpected type: %v", sch3)
				}

				if reflect.DeepEqual(*sch4, sch2) == false {
					t.Errorf("should be equal")
				}
			}()
		})
	}
}

// Rule #2: when received http response code == 200 and not empty response,
// the kafka message value should be changed
func TestHTTPDecoder_200(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	tests := []struct {
		response             string
		expectedMessageValue string
	}{
		// message content updated with new content
		{"response", "response"},
		// message content should not be changed when empty response
		{"", "content"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			server := helper.MockServer(200, tt.response)
			defer server.Close()

			dec := httpdecoder.Decoder{
				URL: server.URL,
			}

			topic := "topic"
			sch := kafka.Schedule{
				Message: &confluent.Message{
					TopicPartition: confluent.TopicPartition{
						Topic: &topic,
					},
					Key:   []byte("video-1"),
					Value: []byte("content"),
				},
			}

			sch2, err := dec.Decode(&sch)
			if err != nil {
				t.Error("unexpected error")
			}

			sch3, ok := sch2.(*kafka.Schedule)
			if !ok {
				t.Errorf("unexpected type: %T", sch2)
			}

			if string(sch3.Message.Value) != tt.expectedMessageValue {
				t.Fatalf("message content not correct: %q", string(sch.Message.Value))
			}
		})
	}
}
