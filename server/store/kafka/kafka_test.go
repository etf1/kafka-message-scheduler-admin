// INTEGRATION TESTS
package kafka_test

import (
	"fmt"
	"testing"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/kafka"
)

// Rule #1: Get by schedule id should return the corresponding schedule(s) : one or many versions
func TestKafkaStore_Get(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics := CreateTopics(t, 2, []int{1, 1}, "schedules")

	kstore, err := kafka.NewStore([]kafka.Bucket{
		{"scheduler-1", GetBootstrapServers(), []string{topics[0]}},
		{"scheduler-2", GetBootstrapServers(), []string{topics[1]}},
	})
	if err != nil {
		t.Fatalf("failed to create kafka store: %v\n", err)
	}
	defer kstore.Close()

	msgs := make([]*confluent.Message, 2000)
	for i := 0; i < 2000; i++ {
		msgs[i] = Message(topics[ifelse(i < 1000, 0, 1)], fmt.Sprintf("schedule-%v", i), "value", now.Add(1*time.Hour).Unix())
	}

	// version 1
	msgs[2] = Message(topics[0], "schedule-1", "value", now.Add(2*time.Hour).Unix())
	// version 2
	msgs[3] = Message(topics[0], "schedule-1", "value", now.Add(3*time.Hour).Unix())

	msgs1 := msgs[:1000]
	msgs2 := msgs[1000:]

	// fmt.Printf("%+v\n", msgs1)
	// fmt.Printf("%+v\n", msgs2)

	ProduceMessages(t, msgs1)
	ProduceMessages(t, msgs2)

	AssertMessagesinTopic(t, topics[0], msgs1)
	AssertMessagesinTopic(t, topics[1], msgs2)

	tests := []struct {
		schedulerName      string
		scheduleID         string
		expectedCount      int
		expectedScheduleID string
	}{
		{"scheduler-1", "schedule-1", 3, "schedule-1"},
		{"scheduler-2", "schedule-1001", 1, "schedule-1001"},
		{"scheduler-2", "schedule-1", 0, ""},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			arr, err := kstore.Get(tt.schedulerName, tt.scheduleID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v := len(arr); v != tt.expectedCount {
				t.Fatalf("unexpected schedules length: %v", v)
			}
			if tt.expectedCount > 0 && arr[0].ID() != tt.expectedScheduleID {
				t.Fatalf("unexpected schedule id: %v", arr[0].ID())
			}
		})
	}
}

// Rule #2: Get by schedule id should return the list of all versions
func TestKafkaStore_Get_many(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics := CreateTopics(t, 1, []int{1}, "schedules")

	kstore, err := kafka.NewStore([]kafka.Bucket{
		{"scheduler-1", GetBootstrapServers(), []string{topics[0]}},
	})
	if err != nil {
		t.Fatalf("failed to create kafka store: %v\n", err)
	}
	defer kstore.Close()

	epoch0 := now.Add(1 * time.Hour).Unix()
	epoch1 := now.Add(2 * time.Hour).Unix()
	epoch2 := now.Add(3 * time.Hour).Unix()

	msgs := []*confluent.Message{
		Message(topics[0], "schedule-1", "value", epoch0),
		Message(topics[0], "schedule-1", "value", epoch1),
		Message(topics[0], "schedule-1", "value", epoch2),
	}

	ProduceMessages(t, msgs)
	AssertMessagesinTopic(t, topics[0], msgs)

	arr, err := kstore.Get("scheduler-1", "schedule-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("arr=%v", arr)
	if v := len(arr); v != 3 {
		t.Fatalf("unexpected schedules length: %v", v)
	}
	if arr[0].Epoch() != epoch0 {
		t.Fatalf("unexpected schedule: %v", arr[0].Epoch())
	}
	if arr[1].Epoch() != epoch1 {
		t.Fatalf("unexpected schedule: %v", arr[0].Epoch())
	}
	if arr[2].Epoch() != epoch2 {
		t.Fatalf("unexpected schedule: %v", arr[0].Epoch())
	}
}

// Rule #2: List should return all schedules by scheduler
func TestKafkaStore_List(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics := CreateTopics(t, 2, []int{1, 1}, "schedules")

	kstore, err := kafka.NewStore([]kafka.Bucket{
		{"scheduler-1", GetBootstrapServers(), []string{topics[0]}},
		{"scheduler-2", GetBootstrapServers(), []string{topics[1]}},
	})
	if err != nil {
		t.Fatalf("failed to create kafka store: %v\n", err)
	}
	defer kstore.Close()

	msgs := make([]*confluent.Message, 2000)
	for i := 0; i < 2000; i++ {
		msgs[i] = Message(topics[ifelse(i < 1000, 0, 1)], fmt.Sprintf("schedule-%v", i), "value", now.Add(1*time.Hour).Unix())
	}

	msgs1 := msgs[:1000]
	msgs2 := msgs[1000:]

	ProduceMessages(t, msgs1)
	ProduceMessages(t, msgs2)

	AssertMessagesinTopic(t, topics[0], msgs1)
	AssertMessagesinTopic(t, topics[1], msgs2)

	tests := []struct {
		schedulerName string
		expectedCount int
	}{
		{"scheduler-1", len(msgs1)},
		{"scheduler-2", len(msgs2)},
		{"scheduler-3", 0},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			lst, err := kstore.List(tt.schedulerName)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			count := 0
			for {
				_, ok := <-lst
				if !ok {
					break
				}
				count++
			}
			if v := count; v != tt.expectedCount {
				t.Fatalf("unexpected schedules length: %v", v)
			}
		})
	}

}

// Rule #3: watch should stream all schedules by type (upsert or deleted)
func TestKafkaStore_Watch(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics := CreateTopics(t, 2, []int{1, 1}, "schedules")

	kstore, err := kafka.NewStore([]kafka.Bucket{
		{"scheduler-1", GetBootstrapServers(), []string{topics[0]}},
		{"scheduler-2", GetBootstrapServers(), []string{topics[1]}},
	})
	if err != nil {
		t.Fatalf("failed to create kafka store: %v\n", err)
	}
	defer kstore.Close()

	msgs := make([]*confluent.Message, 2000)
	for i := 0; i < 2000; i++ {
		var value interface{} = "value"
		if i%2 == 0 {
			value = nil
		}
		msgs[i] = Message(topics[i%2], fmt.Sprintf("schedule-%v", i), value, now.Add(1*time.Hour).Unix())
	}

	msgs1 := msgs[:1000]
	msgs2 := msgs[1000:]

	ProduceMessages(t, msgs1)
	ProduceMessages(t, msgs2)

	lst, err := kstore.Watch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := 0
	upsert := 0
	deleted := 0

loop:
	for {
		select {
		case evt, ok := <-lst:
			if !ok {
				break
			}
			total++
			if evt.EventType == store.DeletedType {
				deleted++
			}
			if evt.EventType == store.UpsertType {
				upsert++
			}
		case <-time.After(2 * time.Second):
			break loop
		}
	}

	if total != 2000 {
		t.Fatalf("unexpected total count: %v", total)
	}
	if upsert != 1000 {
		t.Fatalf("unexpected upsert count: %v", upsert)
	}
	if deleted != 1000 {
		t.Fatalf("unexpected deleted count: %v", deleted)
	}
}
