package pubsub

import "context"

// Producer commnet
type Producer interface {
	Configure(brokers []string, clientID string, topic string) error
	Push(parent context.Context, key, value []byte) (err error)
	Close()
}
