package timings

import (
	"context"
	"fmt"
	"io"
	_ "log"
	"sync/atomic"
	"time"
)

// type CounterMonitor implements the Monitor interface providing a background timings mechanism that tracks incrementing events.
type CounterMonitor struct {
	Monitor
	done_ch chan bool
	start   time.Time
	counter int64
	ticker  *time.Ticker
}

// NewCounterMonitor creates a new Monitor instance that will dispatch notifications using a time.Ticker and 'd'
func NewCounterMonitor(ctx context.Context, d time.Duration) (Monitor, error) {

	done_ch := make(chan bool)
	count := int64(0)

	ticker := time.NewTicker(d)

	t := &CounterMonitor{
		done_ch: done_ch,
		ticker:  ticker,
		counter: count,
	}

	return t, nil
}

// Start will cause background monitoring to begin, dispatching notifications to wr.
func (t *CounterMonitor) Start(ctx context.Context, wr io.Writer) error {

	if !t.start.IsZero() {
		return fmt.Errorf("Monitor has already been started")
	}

	now := time.Now()
	t.start = now

	go func() {

		for {
			select {
			case <-t.done_ch:
				return
			case <-ctx.Done():
				return
			case <-t.ticker.C:
				msg := fmt.Sprintf("processed %d records in %v (started %v)\n", atomic.LoadInt64(&t.counter), time.Since(t.start), t.start)
				wr.Write([]byte(msg))
			}
		}
	}()

	return nil
}

// Stop will cause background monitoring to be halted.
func (t *CounterMonitor) Stop(ctx context.Context) error {
	t.done_ch <- true
	return nil
}

// Signal will cause the background monitors counter to be incremented by one.
func (t *CounterMonitor) Signal(ctx context.Context, args ...interface{}) error {
	return t.increment(ctx)
}

func (t *CounterMonitor) increment(ctx context.Context) error {
	atomic.AddInt64(&t.counter, 1)
	return nil
}
