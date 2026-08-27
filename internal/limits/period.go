package limits

import (
	"paymentconsole/internal/model"
	"sync"
	"time"
)

type PeriodLedger struct {
	mu     sync.Mutex
	day    time.Time
	ledger *Ledger
}

func NewPeriod() *PeriodLedger {
	return &PeriodLedger{day: time.Now().UTC().Truncate(24 * time.Hour), ledger: New()}
}
func (p *PeriodLedger) roll() {
	now := time.Now().UTC().Truncate(24 * time.Hour)
	if now.After(p.day) {
		p.day = now
		p.ledger.Reset()
	}
}
func (p *PeriodLedger) Reserve(c model.Channel, a int64) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.roll()
	return p.ledger.Reserve(c, a)
}
func (p *PeriodLedger) Release(id string, a int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.roll()
	p.ledger.Release(id, a)
}
func (p *PeriodLedger) Used(id string) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.roll()
	return p.ledger.Used(id)
}
func (p *PeriodLedger) Day() time.Time { p.mu.Lock(); defer p.mu.Unlock(); return p.day }
