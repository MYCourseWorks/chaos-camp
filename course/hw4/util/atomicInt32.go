package util

import "sync/atomic"

// AtomicInt32 comment
type AtomicInt32 int32

// Add commnet
func (id *AtomicInt32) Add(other int32) {
	toInt32 := int32(*id)
	atomic.AddInt32(&toInt32, other)
}

// Load commnet
func (id *AtomicInt32) Load() int32 {
	toInt32 := int32(*id)
	return atomic.LoadInt32(&toInt32)
}
