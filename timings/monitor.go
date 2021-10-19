package timings

import (
	"context"
	"fmt"
	"io"
	_ "log"
	"sync/atomic"
	"time"
)

type Monitor struct {
	done_ch chan bool
	start   time.Time
	counter int64
	ticker  *time.Ticker
}

func NewMonitor(ctx context.Context, d time.Duration) (*Monitor, error) {

	done_ch := make(chan bool)
	count := int64(0)

	ticker := time.NewTicker(d)

	t := &Monitor{
		done_ch: done_ch,
		ticker:  ticker,
		counter: count,
	}

	return t, nil
}

func (t *Monitor) Start(ctx context.Context, wr io.Writer) {

	if !t.start.IsZero() {
		return
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

}

func (t *Monitor) Stop(ctx context.Context) {
	t.done_ch <- true
}

func (t *Monitor) Increment(ctx context.Context) {
	atomic.AddInt64(&t.counter, 1)
}
