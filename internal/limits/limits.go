package limits

import (
	"errors"
	"paymentconsole/internal/model"
	"sync"
)

type Ledger struct {
	mu   sync.Mutex
	used map[string]int64
}

func New() *Ledger { return &Ledger{used: map[string]int64{}} }
func (l *Ledger) Reserve(c model.Channel, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	next := l.used[c.ID] + amount
	if c.LimitCents > 0 && next > c.LimitCents {
		return errors.New("daily limit exceeded")
	}
	l.used[c.ID] = next
	return nil
}
func (l *Ledger) Release(id string, amount int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.used[id] -= amount
	if l.used[id] < 0 {
		l.used[id] = 0
	}
}
func (l *Ledger) Used(id string) int64 { l.mu.Lock(); defer l.mu.Unlock(); return l.used[id] }
func (l *Ledger) Reset()               { l.mu.Lock(); l.used = map[string]int64{}; l.mu.Unlock() }
func (l *Ledger) Snapshot() map[string]int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	m := map[string]int64{}
	for k, v := range l.used {
		m[k] = v
	}
	return m
}
func (l *Ledger) CanReserve(c model.Channel, a int64) bool { return l.Reserve(c, a) == nil }
