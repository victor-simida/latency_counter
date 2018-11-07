package latency

import (
	"sync/atomic"
	"time"
)

// if a latency haven't got new input in last 300 seconds, set it with zero
var timeout int64 = 300

// Latency ...
type Latency struct {
	record         []int32
	index          int8
	averageLatency int32
	input          chan byte
	lastInsertTime int64
}

// NewLatency ...
func NewLatency() *Latency {
	return &Latency{
		record: make([]int32, 10),
		input:  make(chan byte, 1),
	}
}

// Input ...
func (l *Latency) Input(input time.Duration) {
	t := int32(input / time.Millisecond)
	if t <= 0 {
		return
	}

	select {
	case l.input <- 1:
	default:
		return
	}

	// If last access time is long before, reset the record
	if time.Now().Unix()-atomic.LoadInt64(&l.lastInsertTime) >= timeout {
		l.index = 0
		for i := 0; i < len(l.record); i++ {
			l.record[i] = 0
		}
	}

	old := l.record[l.index]
	// If record is not filled, insert directly
	if old == 0 {
		l.record[l.index] = t
		atomic.StoreInt32(&l.averageLatency, (l.averageLatency*int32(l.index)+t)/(int32(l.index)+1))
		l.index = l.index + 1
		if l.index == 10 {
			l.index = 0
		}
	} else {
		l.record[l.index] = t
		l.index = l.index + 1
		if l.index == 10 {
			l.index = 0
		}
		atomic.StoreInt32(&l.averageLatency, l.averageLatency+(t-old)/(10))
	}
	atomic.StoreInt64(&l.lastInsertTime, time.Now().Unix())
	<-l.input
}

// Get ...
func (l *Latency) Get() int32 {
	if time.Now().Unix()-atomic.LoadInt64(&l.lastInsertTime) >= timeout {
		return 0
	}
	return atomic.LoadInt32(&l.averageLatency)
}
