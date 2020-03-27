package pubsub

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MartinNikolovMarinov/sb-infra/infra"
	"github.com/segmentio/kafka-go"

	// Need this to read the messages
	_ "github.com/segmentio/kafka-go/snappy"
)

// kafkaCache commnet
type kafkaCache struct {
	reader *kafka.Reader
	data   map[string]CacheMsg
}

// NewKafkaCache comment
func NewKafkaCache(brokers []string, topic string, partition int) (Cache, error) {
	var err error
	cache := &kafkaCache{}

	if brokers == nil {
		localIP := os.Getenv("MY_IP")
		if localIP == "" {
			return nil, fmt.Errorf("Failed to create cache for topic: %s", topic)
		}
		brokers = []string{localIP + ":9092"}
	}

	// make a new reader that consumes from topic-A
	config := kafka.ReaderConfig{
		Brokers:         brokers,
		Topic:           topic,
		Partition:       partition,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         4 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}

	reader := kafka.NewReader(config)
	err = reader.SetOffset(0)
	if err != nil {
		return nil, fmt.Errorf("Failed to set cache offset for topic: %s", topic)
	}

	cache.reader = reader
	cache.data = make(map[string]CacheMsg)
	go cache.startReading()
	return cache, nil
}

func (c *kafkaCache) startReading() {
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			infra.Error(err.Error())
			continue
		}

		msg := CacheMsg{
			Key:       string(m.Key),
			Topic:     m.Topic,
			Partition: m.Partition,
			Offset:    m.Offset,
			Value:     m.Value,
			Timestamp: m.Time,
		}

		if msg.Key == "" {
			infra.Error("Invalid key from topic: %s", msg.Topic)
			continue
		}

		c.data[msg.Key] = msg
	}
}

// Get commnet
func (c *kafkaCache) Get(id string) *CacheMsg {
	if v, ok := c.data[id]; ok {
		return &v
	}

	return nil
}

// Len commnet
func (c *kafkaCache) Len() int {
	return len(c.data)
}

// All commnet
func (c *kafkaCache) All() []CacheMsg {
	i := 0
	ret := make([]CacheMsg, len(c.data))
	for _, v := range c.data {
		ret[i] = v
		i++
	}

	return ret
}

// Close comment
func (c *kafkaCache) Close() {
	c.reader.Close()
}
