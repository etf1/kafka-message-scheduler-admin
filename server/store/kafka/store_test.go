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

var (
	ifelse = helper.IfElse
)

// TODO : use helper package for ProduceMessages etc ...

// Rule #1: Get by schedule id should return the corresponding schedule(s) : one or many versions
func TestKafkaStore_Get(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics, err := helper.CreateTopics(2, []int{1, 1}, "schedules")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	kstore, err := kafka.NewStore([]kafka.Bucket{
		{"scheduler-1", helper.GetDefaultBootstrapServers(), []string{topics[0]}},
		{"scheduler-2", helper.GetDefaultBootstrapServers(), []string{topics[1]}},
	})
	if err != nil {
		t.Errorf("failed to create kafka store: %v\n", err)
	}
	defer kstore.Close()

	msgs := make([]*confluent.Message, 20)
	for i := 0; i < 20; i++ {
		msgs[i] = helper.Message(topics[ifelse(i < 10, 0, 1)], fmt.Sprintf("schedule-%v", i), "value", now.Add(1*time.Hour).Unix())
	}

	msgs1 := msgs[:10]
	msgs2 := msgs[10:]

	helper.ProduceMessages(msgs1)
	helper.ProduceMessages(msgs2)

	err = helper.AssertMessagesinTopic(topics[0], msgs1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	err = helper.AssertMessagesinTopic(topics[1], msgs2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tests := []struct {
		schedulerName      string
		scheduleID         string
		expectedCount      int
		expectedScheduleID string
	}{
		{"scheduler-1", "schedule-1", 1, "schedule-1"},
		{"scheduler-1", "schedule-11", 0, ""},

		{"scheduler-2", "schedule-1", 0, ""},
		{"scheduler-2", "schedule-11", 1, "schedule-11"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v", i+1), func(t *testing.T) {
			arr, err := kstore.Get(tt.schedulerName, tt.scheduleID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if v := len(arr); v != tt.expectedCount {
				t.Errorf("unexpected schedules length: %v %+v", v, arr)
			}
			if tt.expectedCount > 0 && arr[0].ID() != tt.expectedScheduleID {
				t.Errorf("unexpected schedule id: %v", arr[0].ID())
			}
		})
	}
}

// Rule #2: Get by schedule id should return the list of all versions
func TestKafkaStore_Get_many(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics, err := helper.CreateTopics(1, []int{1}, "schedules")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	kstore, err := kafka.NewStore([]kafka.Bucket{
		{"scheduler-1", helper.GetDefaultBootstrapServers(), []string{topics[0]}},
	})
	if err != nil {
		t.Errorf("failed to create kafka store: %v\n", err)
	}
	defer kstore.Close()

	epoch0 := now.Add(1 * time.Hour).Unix()
	epoch1 := now.Add(2 * time.Hour).Unix()
	epoch2 := now.Add(3 * time.Hour).Unix()

	msgs := []*confluent.Message{
		helper.Message(topics[0], "schedule-1", "value", epoch0),
		helper.Message(topics[0], "schedule-1", "value", epoch1),
		helper.Message(topics[0], "schedule-1", "value", epoch2),
	}

	helper.ProduceMessages(msgs)
	err = helper.AssertMessagesinTopic(topics[0], msgs)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	arr, err := kstore.Get("scheduler-1", "schedule-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if v := len(arr); v != 3 {
		t.Errorf("unexpected schedules length: %v", v)
	}
	if arr[0].Epoch() != epoch0 {
		t.Errorf("unexpected schedule: %v", arr[0].Epoch())
	}
	if arr[1].Epoch() != epoch1 {
		t.Errorf("unexpected schedule: %v", arr[0].Epoch())
	}
	if arr[2].Epoch() != epoch2 {
		t.Errorf("unexpected schedule: %v", arr[0].Epoch())
	}
}

// Rule #2: List should return all schedules by scheduler
func TestKafkaStore_List(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics, err := helper.CreateTopics(2, []int{1, 1}, "schedules")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	kstore, err := kafka.NewStore([]kafka.Bucket{
		{"scheduler-1", helper.GetDefaultBootstrapServers(), []string{topics[0]}},
		{"scheduler-2", helper.GetDefaultBootstrapServers(), []string{topics[1]}},
	})
	if err != nil {
		t.Errorf("failed to create kafka store: %v\n", err)
	}
	defer kstore.Close()

	msgs := make([]*confluent.Message, 20)
	for i := 0; i < 20; i++ {
		msgs[i] = helper.Message(topics[ifelse(i < 10, 0, 1)], fmt.Sprintf("schedule-%v", i), "value", now.Add(1*time.Hour).Unix())
	}

	msgs1 := msgs[:10]
	msgs2 := msgs[10:]

	helper.ProduceMessages(msgs1)
	helper.ProduceMessages(msgs2)

	err = helper.AssertMessagesinTopic(topics[0], msgs1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = helper.AssertMessagesinTopic(topics[1], msgs2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

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
				t.Errorf("unexpected error: %v", err)
				return
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
				t.Errorf("unexpected schedules length: %v", v)
			}
		})
	}
}

// Rule #3: watch should stream all schedules by type (upsert or deleted)
func TestKafkaStore_Watch(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics, err := helper.CreateTopics(2, []int{1, 1}, "schedules")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	kstore, err := kafka.NewStore([]kafka.Bucket{
		{"scheduler-1", helper.GetDefaultBootstrapServers(), []string{topics[0]}},
		{"scheduler-2", helper.GetDefaultBootstrapServers(), []string{topics[1]}},
	})
	if err != nil {
		t.Errorf("failed to create kafka store: %v\n", err)
	}
	defer kstore.Close()

	msgs := make([]*confluent.Message, 20)
	for i := 0; i < 20; i++ {
		var value interface{} = "value"
		if i%2 == 0 {
			value = nil
		}
		topic := topics[0]
		if i >= 10 {
			topic = topics[1]
		}
		msgs[i] = helper.Message(topic, fmt.Sprintf("schedule-%v", i), value, now.Add(1*time.Hour).Unix())
	}

	msgs1 := msgs[:10]
	msgs2 := msgs[10:]

	helper.ProduceMessages(msgs1)
	helper.ProduceMessages(msgs2)

	err = helper.AssertMessagesinTopic(topics[0], msgs1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = helper.AssertMessagesinTopic(topics[1], msgs2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	lst, err := kstore.Watch()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
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

	t.Logf("upsert=%v deleted=%v total=%v", upsert, deleted, total)

	if total != 20 {
		t.Errorf("unexpected total count: %v", total)
	}
	if upsert != 10 {
		t.Errorf("unexpected upsert count: %v", upsert)
	}
	if deleted != 10 {
		t.Errorf("unexpected deleted count: %v", deleted)
	}
}
