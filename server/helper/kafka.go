package helper

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/etf1/kafka-message-scheduler/schedule"
	kafka_schedule "github.com/etf1/kafka-message-scheduler/schedule/kafka"
	"github.com/etf1/kafka-message-scheduler/schedule/simple"
	log "github.com/sirupsen/logrus"
)

// This package contains reusable code for tests.

const (
	AdminTimeout = 2 * time.Second
	ReadTimeout  = 5 * time.Second
	FlushTimeout = 10000
)

// Converts int, string to slice of byte
func toBytes(value interface{}) []byte {
	if value == nil {
		return nil
	}

	k := "unknow type"
	switch kt := value.(type) {
	case int:
		k = strconv.Itoa(kt)
	case string:
		k = kt
	case []byte:
		return kt
	}

	return []byte(k)
}

func NewKafkaSchedule(topic string, key, value interface{}, epoch int64, targetTopic string, targetKey interface{}) kafka_schedule.Schedule {
	headers := []confluent.Header{
		{
			Key:   kafka_schedule.Epoch,
			Value: toBytes(strconv.FormatInt(epoch, 10)),
		},
		{
			Key:   kafka_schedule.TargetTopic,
			Value: toBytes(targetTopic),
		},
		{
			Key:   kafka_schedule.TargetKey,
			Value: toBytes(targetKey),
		}}

	msg := &confluent.Message{
		TopicPartition: confluent.TopicPartition{Topic: &topic, Partition: confluent.PartitionAny},
		Headers:        headers,
		Value:          toBytes(value),
		Key:            toBytes(key),
		Timestamp:      time.Now(),
	}

	return kafka_schedule.Schedule{Message: msg}
}

// Creates a kafka.Message
func Message(topic string, key, value interface{}, epoch int64) *confluent.Message {
	headers := []confluent.Header{
		{
			Key:   kafka_schedule.Epoch,
			Value: []byte(strconv.FormatInt(epoch, 10)),
		}}

	return &confluent.Message{
		TopicPartition: confluent.TopicPartition{Topic: &topic, Partition: confluent.PartitionAny},
		Headers:        headers,
		Value:          toBytes(value),
		Key:            toBytes(key),
	}
}

// Creates multiple simple.Schedules
func SimpleSchedules(count int) []schedule.Schedule {
	result := make([]schedule.Schedule, 0)
	now := time.Now()
	for i := 0; i < count; i++ {
		epoch := now.Add(time.Duration(i) * time.Second)
		ts := now.Add(time.Duration(i) * time.Second)
		result = append(result, simple.NewSchedule(fmt.Sprintf("schedule-%v", i), epoch, ts))
	}
	return result
}

// tells if the tests is running in docker
// why ? because the kafka host and port are different in docker than outside docker
func isRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); os.IsNotExist(err) {
		return false
	}
	return true
}

// Get the bootstrap servers because in or out the docker the kafka server is different
func GetDefaultBootstrapServers() string {
	if isRunningInDocker() {
		fmt.Println("kafka bootstrap servers=kafka:29092")
		return "kafka:29092"
	}
	fmt.Println("kafka bootstrap servers=localhost:9092")
	return "localhost:9092"
}

// creates a random topic name, each test got a different topic name
func RandomTopicName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, RandNum())
}

// // Creates topics based on the array of nbPartitions param
func CreateTopics(nbTopic int, nbPartitions []int, prefix string) ([]string, error) {
	topics := make([]string, nbTopic)

	adm, err := confluent.NewAdminClient(&confluent.ConfigMap{
		"bootstrap.servers": GetDefaultBootstrapServers(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create admin client: %v", err)
	}
	defer adm.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	replicationFactor := 1
	specs := make([]confluent.TopicSpecification, nbTopic)
	for i := 0; i < len(specs); i++ {
		topics[i] = RandomTopicName(prefix)
		specs[i] = confluent.TopicSpecification{
			Topic:             topics[i],
			NumPartitions:     nbPartitions[i],
			ReplicationFactor: replicationFactor}
	}

	results, err := adm.CreateTopics(ctx, specs, confluent.SetAdminOperationTimeout(AdminTimeout))
	if err != nil {
		return nil, fmt.Errorf("failed to create topics %v: %v", topics, err)
	}

	for _, result := range results {
		log.Printf("%s", result)
	}

	return topics, nil
}

func ReCreateTopic(topic string) error {
	adm, err := confluent.NewAdminClient(&confluent.ConfigMap{
		"bootstrap.servers": GetDefaultBootstrapServers(),
	})
	if err != nil {
		return fmt.Errorf("failed to create admin client: %w", err)
	}
	defer adm.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results, err := adm.DeleteTopics(ctx, []string{topic}, confluent.SetAdminOperationTimeout(AdminTimeout))
	if err != nil {
		return fmt.Errorf("failed to delete topic %v: %w", topic, err)
	}

	sleepDuration := 2 * time.Second
	time.Sleep(sleepDuration)

	for _, result := range results {
		log.Debugf("result deletion: %s\n", result)
	}
	spec := confluent.TopicSpecification{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	results, err = adm.CreateTopics(ctx, []confluent.TopicSpecification{spec}, confluent.SetAdminOperationTimeout(AdminTimeout))
	if err != nil {
		return fmt.Errorf("failed to create topic %v: %w", topic, err)
	}

	time.Sleep(sleepDuration)

	for _, result := range results {
		log.Debugf("result creation: %s\n", result)
	}

	return nil
}

func GetDefaultProducerConfig() *confluent.ConfigMap {
	return &confluent.ConfigMap{
		"bootstrap.servers": GetDefaultBootstrapServers(),
	}
}

func ProduceMessages(messages []*confluent.Message) error {
	p, err := confluent.NewProducer(GetDefaultProducerConfig())
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	defer p.Close()

	deliveryChan := make(chan confluent.Event)
	defer close(deliveryChan)

	for _, message := range messages {
		err := p.Produce(message, deliveryChan)
		if err != nil {
			return fmt.Errorf("produce failed: %w", err)
		}
		e := <-deliveryChan
		m := e.(*confluent.Message)

		if m.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
		}

		log.Debugf("delivered message %v to topic %s [%d] at offset %v\n",
			string(m.Key), *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	p.Flush(FlushTimeout)

	return nil
}

func GetDefaultConsumerConfig(prefix string) *confluent.ConfigMap {
	return &confluent.ConfigMap{
		"bootstrap.servers":     GetDefaultBootstrapServers(),
		"broker.address.family": "v4",
		"group.id":              fmt.Sprintf("%s-%d", prefix, RandNum()),
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	}
}

// Compares 2 kafka messages
func AssertMessageEquals(m1, m2 *confluent.Message) bool {
	if !reflect.DeepEqual(m1.Headers, m2.Headers) {
		log.Printf("headers not equals, %v != %v", m1.Headers, m2.Headers)
		return false
	}
	if !bytes.Equal(m1.Key, m2.Key) {
		log.Printf("keys not equals, %v != %v", string(m1.Key), string(m2.Key))
		return false
	}
	if !bytes.Equal(m1.Value, m2.Value) {
		log.Printf("values not equals, %v != %v", string(m1.Value), string(m2.Value))
		return false
	}

	return true
}

// Verifies if specified messages are in the topic
func AssertMessagesinTopic(topic string, msgs []*confluent.Message) error {
	config := GetDefaultConsumerConfig("cg-test")
	log.Printf("consumer config: topic=%v %+v", topic, config)

	c, err := confluent.NewConsumer(config)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %v", err)
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}

	foundCount := 0
	for _, m := range msgs {
		msg, err := c.ReadMessage(ReadTimeout)
		if err != nil {
			if err.(confluent.Error).Code() == confluent.ErrTimedOut {
				break
			}
			return fmt.Errorf("unexpected error: %v", err)
		}
		AssertMessageEquals(m, msg)
		foundCount++
	}
	if foundCount != len(msgs) {
		return fmt.Errorf("unexpected found count %v, expected %v", foundCount, len(msgs))
	}

	return nil
}

func CopyKafkaSchedule(s kafka_schedule.Schedule) kafka_schedule.Schedule {
	s2 := s
	msg2 := *s.Message
	s2.Message = &msg2

	return s2
}

type KafkaMessageSimpleDecoder struct {
	Called int
}

func (k *KafkaMessageSimpleDecoder) Decode(s schedule.Schedule) (schedule.Schedule, error) {
	k.Called++
	return s, nil
}
