package pubsub

import "context"

// MetricsProducer comment
type MetricsProducer struct {
	p Producer
}

const metrixTopic = "Metrics"

// NewMetrics commnet
func NewMetrics() (*MetricsProducer, error) {
	p := NewKafkaProducer()
	err := p.Configure(nil, "consumer", metrixTopic)
	return &MetricsProducer{p: p}, err
}

// Publish comment
func (m *MetricsProducer) Publish(id []byte, data []byte) error {
	err := m.p.Push(context.Background(), id, data)
	if err != nil {
		return err
	}

	return nil
}

// Close comment
func (m *MetricsProducer) Close() {
	m.p.Close()
}
