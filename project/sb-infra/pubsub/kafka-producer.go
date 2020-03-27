package pubsub

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

// kafkaProducer commnet
type kafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer comment
func NewKafkaProducer() Producer {
	return &kafkaProducer{}
}

// Configure comment
func (p *kafkaProducer) Configure(brokers []string, clientID string, topic string) error {
	if p == nil {
		p = &kafkaProducer{}
	}

	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientID,
	}

	if brokers == nil {
		localIP := os.Getenv("MY_IP")
		if localIP == "" {
			return errors.New("brokers argument is empty")
		}
		b1 := localIP + ":9092"
		brokers = []string{b1}
	}

	config := kafka.WriterConfig{
		Brokers:          brokers,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}

	p.writer = kafka.NewWriter(config)
	return nil
}

// Push comment
func (p *kafkaProducer) Push(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(parent, time.Second*10)
	defer cancel()
	return p.writer.WriteMessages(ctx, message)
}

// Close comment
func (p *kafkaProducer) Close() {
	p.writer.Close()
}
