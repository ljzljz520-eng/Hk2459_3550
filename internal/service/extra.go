package service

import (
	"context"
	"errors"
	"paymentconsole/internal/model"
	"paymentconsole/internal/store"
	"time"
)

func (s *Service) EnsureDefaults(ctx context.Context) error {
	for _, c := range []model.Channel{model.NewChannel("card", model.Card), model.NewChannel("wallet", model.Wallet), model.NewChannel("installment", model.Installment), model.NewChannel("offline", model.Offline)} {
		if _, e := s.store.GetChannel(c.ID); e != nil {
			if e := s.CreateChannel(ctx, c); e != nil {
				return e
			}
		}
	}
	return nil
}
func (s *Service) SetNotice(ctx context.Context, id, notice string) error {
	c, e := s.store.GetChannel(id)
	if e != nil {
		return e
	}
	return s.UpdateChannel(ctx, id, c.Enabled, c.LimitCents, c.FeeBps, notice, "notice")
}
func (s *Service) SetLimit(ctx context.Context, id string, limit int64) error {
	if limit < 0 {
		return errors.New("negative limit")
	}
	c, e := s.store.GetChannel(id)
	if e != nil {
		return e
	}
	return s.UpdateChannel(ctx, id, c.Enabled, limit, c.FeeBps, c.Notice, "limit")
}
func (s *Service) SetFee(ctx context.Context, id string, fee int) error {
	if fee < 0 || fee > 10000 {
		return errors.New("invalid fee")
	}
	c, e := s.store.GetChannel(id)
	if e != nil {
		return e
	}
	return s.UpdateChannel(ctx, id, c.Enabled, c.LimitCents, fee, c.Notice, "fee")
}
func (s *Service) Snapshot() ([]model.Channel, error) { return s.store.ListChannels() }
func (s *Service) CreateWorkflow(name string, steps []string) model.Workflow {
	return model.Workflow{ID: store.ID("wf"), Name: name, Steps: steps, Status: "draft", CreatedAt: time.Now()}
}
