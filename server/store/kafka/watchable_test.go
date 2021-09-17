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

// Rule #1: watch should stream all schedules by type (upsert or deleted)
func TestKafkaWatchableStore_Watch(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics, err := helper.CreateTopics(2, []int{1, 1}, "schedules")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	buckets := []kafka.Bucket{
		{"scheduler-1", helper.GetDefaultBootstrapServers(), []string{topics[0]}},
		{"scheduler-2", helper.GetDefaultBootstrapServers(), []string{topics[1]}},
	}
	kstore, err := kafka.NewWatchableStore(nil, buckets...)
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
		case <-time.After(5 * time.Second):
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

// Rule #2: test AddBuckets, adding bucket after instanciation should work properly
// updating bucket's wrong config or adding new bucket, an event StoreResetType should be thrown
func TestKafkaWatchableStore_AddBuckets(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics, err := helper.CreateTopics(2, []int{1, 1}, "schedules")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// buckets to add
	bucket1 := kafka.Bucket{"scheduler-1", helper.GetDefaultBootstrapServers(), []string{topics[0]}}
	bucket2 := kafka.Bucket{"scheduler-2", "unknown:9092", []string{topics[1]}}
	// fix the bucket config
	bucket2Fix := kafka.Bucket{"scheduler-2", helper.GetDefaultBootstrapServers(), []string{topics[1]}}

	kstore, err := kafka.NewWatchableStore(nil)
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
		msgs[i] = helper.Message(topics[ifelse(i < 10, 0, 1)], fmt.Sprintf("schedule-%v", i), value, now.Add(1*time.Hour).Unix())
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

	go func() {
		kstore.AddBuckets(bucket1, bucket2)
		time.Sleep(2 * time.Second)
		// should fix the scheduler-2 config and process assigned messages
		kstore.AddBuckets(bucket2Fix)
	}()

	total := 0
	upsert := 0
	deleted := 0
	reset := 0

loop:
	for {
		select {
		case evt, ok := <-lst:
			if !ok {
				t.Logf("channel closed, exiting")
				break
			}
			total++
			if evt.EventType == store.DeletedType {
				deleted++
			}
			if evt.EventType == store.UpsertType {
				upsert++
			}
			if evt.EventType == store.StoreResetType {
				reset++
			}
		case <-time.After(10 * time.Second):
			t.Logf("timeout, exiting")
			break loop
		}
	}

	t.Logf("upsert=%v deleted=%v reset=%v total=%v", upsert, deleted, reset, total)

	if total != 21 {
		t.Errorf("unexpected total count: %v", total)
	}
	if reset != 1 {
		t.Errorf("unexpected reset count: %v", reset)
	}
	if upsert != 10 {
		t.Errorf("unexpected upsert count: %v", upsert)
	}
	if deleted != 10 {
		t.Errorf("unexpected deleted count: %v", deleted)
	}
}

// Rule #3: when a decoder is specified, it should be used to update message value (body)
func TestKafkaWatchableStore_Decoder(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	topics, err := helper.CreateTopics(2, []int{1, 1}, "schedules")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	buckets := []kafka.Bucket{
		{"scheduler-1", helper.GetDefaultBootstrapServers(), []string{topics[0]}},
		{"scheduler-2", helper.GetDefaultBootstrapServers(), []string{topics[1]}},
	}

	dec := &helper.KafkaMessageSimpleDecoder{}
	kstore, err := kafka.NewWatchableStore(dec, buckets...)
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
		case <-time.After(5 * time.Second):
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
	if dec.Called != 20 {
		t.Errorf("unexpected decoder count: %v", dec.Called)
	}
}
