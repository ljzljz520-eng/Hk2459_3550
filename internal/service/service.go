package service

import (
	"context"
	"errors"
	"paymentconsole/internal/model"
	"paymentconsole/internal/store"
	"sync"
	"time"
)

var ErrUnavailable = errors.New("payment channel unavailable")
var ErrLimit = errors.New("payment exceeds channel limit")

type Service struct {
	store    *store.Store
	mu       sync.RWMutex
	cache    map[string]model.Channel
	inflight map[string]bool
}

func New(s *store.Store) *Service {
	return &Service{store: s, cache: map[string]model.Channel{}, inflight: map[string]bool{}}
}
func (s *Service) load(ctx context.Context, id string) (model.Channel, error) {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return model.Channel{}, ctx.Err()
		default:
		}
	}
	s.mu.RLock()
	c, ok := s.cache[id]
	s.mu.RUnlock()
	if ok {
		return c, nil
	}
	c, e := s.store.GetChannel(id)
	if e == nil {
		s.mu.Lock()
		s.cache[id] = c
		s.mu.Unlock()
	}
	return c, e
}
func (s *Service) CreateChannel(ctx context.Context, c model.Channel) error {
	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	if c.ID == "" || c.Name == "" {
		return errors.New("invalid channel")
	}
	c.Version = 1
	c.UpdatedAt = time.Now()
	if e := s.store.SaveChannel(c); e != nil {
		return e
	}
	s.mu.Lock()
	s.cache[c.ID] = c
	s.mu.Unlock()
	return nil
}
func (s *Service) UpdateChannel(ctx context.Context, id string, enabled bool, limit int64, fee int, notice, actor string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	old, e := s.store.GetChannel(id)
	if e != nil {
		return e
	}
	next := old
	next.Enabled = enabled
	next.LimitCents = limit
	next.FeeBps = fee
	next.Notice = notice
	next.Version++
	next.UpdatedAt = time.Now()
	if e = s.store.SaveChannel(next); e != nil {
		return e
	}
	a := model.AuditEvent{ID: store.ID("audit"), Action: "channel_update", ChannelID: id, Before: old.Enabled, After: enabled, Actor: actor, CreatedAt: time.Now()}
	if e = s.store.SaveAudit(a); e != nil {
		return e
	}
	// Always refresh the in-memory cache with the persisted state. The cache is
	// read-through on every Pay request, so a stale entry here would let new
	// payments be routed to a channel that was just switched off.
	s.mu.Lock()
	s.cache[id] = next
	s.mu.Unlock()
	return nil
}
func (s *Service) Preview(ctx context.Context, id string, enabled bool) (string, error) {
	c, e := s.load(ctx, id)
	if e != nil {
		return "", e
	}
	if c.Enabled == enabled {
		return "no impact", nil
	}
	if enabled {
		return "traffic will resume", nil
	}
	return "new payments will be rejected", nil
}
func (s *Service) Pay(ctx context.Context, r model.PaymentRequest) (model.PaymentResult, error) {
	c, e := s.load(ctx, r.ChannelID)
	if e != nil {
		return model.PaymentResult{ID: r.ID, ChannelID: r.ChannelID, Reason: ErrUnavailable.Error()}, ErrUnavailable
	}
	if !c.CanAccept(r.AmountCents) {
		if !c.Enabled {
			return model.PaymentResult{ID: r.ID, ChannelID: r.ChannelID, Reason: ErrUnavailable.Error()}, ErrUnavailable
		}
		return model.PaymentResult{ID: r.ID, ChannelID: r.ChannelID, Reason: ErrLimit.Error()}, ErrLimit
	}
	return model.PaymentResult{ID: r.ID, ChannelID: r.ChannelID, Accepted: true, FeeCents: c.Fee(r.AmountCents)}, nil
}
func (s *Service) Invalidate(id string) { s.mu.Lock(); delete(s.cache, id); s.mu.Unlock() }
func (s *Service) Refresh(id string) error {
	c, e := s.store.GetChannel(id)
	if e != nil {
		return e
	}
	s.mu.Lock()
	s.cache[id] = c
	s.mu.Unlock()
	return nil
}
func (s *Service) Close() error { return s.store.Close() }
func (s *Service) BatchUpdate(ctx context.Context, ids []string, enabled bool) map[string]error {
	out := map[string]error{}
	for _, id := range ids {
		if e := s.UpdateChannel(ctx, id, enabled, 0, 0, "", "batch"); e != nil {
			out[id] = e
		}
	}
	return out
}
