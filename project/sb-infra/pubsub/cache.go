package pubsub

import "time"

// CacheMsg comment
type CacheMsg struct {
	Key       string
	Topic     string
	Partition int
	Offset    int64
	Timestamp time.Time
	Value     []byte
}

// Cache comment
type Cache interface {
	Get(id string) *CacheMsg
	Len() int
	All() []CacheMsg
	Close()
}
