package latency

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// NewLatency ...
func NewLatencyMock() *Latency {
	return &Latency{
		averageLatency: 100,
	}
}

func TestLatency(t *testing.T) {
	l := NewLatency()

	for i := 1; i <= 10000; i++ {
		go func() {
			l.Input(time.Second)
		}()
	}

	for i := 1; i <= 10000; i++ {
		assert.Equal(t, int32(time.Second/time.Millisecond), l.Get())
	}
}

func TestLatency2(t *testing.T) {
	l := NewLatency()

	for i := 1; i <= 10000; i++ {
		l.Input(time.Duration(rand.Intn(5)) * time.Second)
	}

	fmt.Println(l.Get())
	for i := 1; i <= 10000; i++ {
		assert.NotEqual(t, int32(0), l.Get())
	}
}

func TestLatency3(t *testing.T) {
	l := NewLatency()

	for i := 1; i <= 10000; i++ {
		l.Input(time.Second)
	}
	for i := 1; i <= 10000; i++ {
		l.Input(time.Millisecond * 500)
	}

	for i := 1; i <= 10000; i++ {
		assert.Equal(t, int32(500), l.Get())
	}
}

func TestLatencyClean(t *testing.T) {
	timeout = 5
	l := NewLatency()
	time.Sleep(time.Second)
	l.Input(time.Second)
	time.Sleep(10 * time.Second)

	assert.Equal(t, int32(0), l.Get())
	timeout = 300
}

func BenchmarkLatencyInput(b *testing.B) {
	l := NewLatency()
	for i := 0; i <= b.N; i++ {
		l.Input(time.Second)
	}
}

func BenchmarkLatencyGet(b *testing.B) {
	l := NewLatency()
	for i := 1; i <= 10000; i++ {
		go func() {
			l.Input(time.Second)
		}()
	}
	for i := 0; i <= b.N; i++ {
		l.Get()
	}
}
func BenchmarkLatencyNew(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		NewLatency()
	}
}

func BenchmarkAtomic(b *testing.B) {
	var temp int32
	for i := 0; i <= b.N; i++ {
		atomic.AddInt32(&temp, 1)
	}
}
func BenchmarkAtomic2(b *testing.B) {
	var temp int32
	for i := 0; i <= b.N; i++ {
		atomic.StoreInt32(&temp, 1)
	}
}
