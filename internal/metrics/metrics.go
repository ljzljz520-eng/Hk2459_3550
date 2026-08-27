package metrics

import (
	"sync"
	"time"
)

type Counter struct {
	mu     sync.Mutex
	values map[string]int64
}

func New() *Counter                      { return &Counter{values: map[string]int64{}} }
func (c *Counter) Inc(k string)          { c.mu.Lock(); c.values[k]++; c.mu.Unlock() }
func (c *Counter) Add(k string, n int64) { c.mu.Lock(); c.values[k] += n; c.mu.Unlock() }
func (c *Counter) Get(k string) int64    { c.mu.Lock(); defer c.mu.Unlock(); return c.values[k] }
func (c *Counter) Snapshot() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := map[string]int64{}
	for k, v := range c.values {
		out[k] = v
	}
	return out
}
func (c *Counter) MarkLatency(start time.Time) { c.Add("latency_ms", time.Since(start).Milliseconds()) }
