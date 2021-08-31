// INTEGRATION TESTS
package kafka_test

import (
	"fmt"
	"testing"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
	"github.com/etf1/kafka-message-scheduler-admin/server/runner/kafka"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
)

// Rule #1: watch should stream all schedules by type (upsert or deleted) (from resolver)
func TestKafkaWatchableStoreFromResolver_Watch(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	now := time.Now()

	err := helper.ReCreateTopic("schedules")
	if err != nil {
		t.Error(err)
		return
	}

	resolver := httpresolver.NewResolver(config.SchedulersAddr())

	kstore, err := kafka.NewWatchableStoreFromResolver(resolver, kafka.DefaultTopics)
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
		msgs[i] = helper.Message("schedules", fmt.Sprintf("schedule-%v", i), value, now.Add(1*time.Hour).Unix())
	}

	msgs1 := msgs[:10]
	msgs2 := msgs[10:]

	helper.ProduceMessages(msgs1)
	helper.ProduceMessages(msgs2)

	err = helper.AssertMessagesinTopic("schedules", msgs)
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
