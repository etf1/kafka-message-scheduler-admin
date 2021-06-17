// INTEGRATION TESTS
package kafka_test

import (
	"fmt"
	"testing"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/kafka"
)

// Rule #1: watch should stream all schedules by type (upsert or deleted)
func TestKafkaWatchableStore_Watch(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics := CreateTopics(t, 2, []int{1, 1}, "schedules")

	kstore, err := kafka.NewWatchableStore([]kafka.Bucket{
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

// Rule #1: watch should stream all schedules by type (upsert or deleted) (from resolver)
func TestKafkaWatchableStoreFromResolver_Watch(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	CreateTopic(t, "scheduler")

	resolver := httpresolver.NewResolver([]string{"localhost"})

	kstore, err := kafka.NewWatchableStoreFromResolver(resolver)
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
		msgs[i] = Message("scheduler", fmt.Sprintf("schedule-%v", i), value, now.Add(1*time.Hour).Unix())
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
