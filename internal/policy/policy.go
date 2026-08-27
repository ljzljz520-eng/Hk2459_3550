package policy

import (
	"errors"
	"paymentconsole/internal/model"
)

type Policy struct {
	RequirePreview bool
	MaxFeeBps      int
	Restricted     map[model.ChannelType]bool
}

func New() Policy {
	return Policy{RequirePreview: true, MaxFeeBps: 1000, Restricted: map[model.ChannelType]bool{}}
}
func (p Policy) Check(c model.Channel) error {
	if c.FeeBps > p.MaxFeeBps {
		return errors.New("fee exceeds policy")
	}
	if p.Restricted[c.Type] && c.Enabled {
		return errors.New("channel restricted")
	}
	return nil
}
func (p *Policy) Restrict(t model.ChannelType, v bool) { p.Restricted[t] = v }
func (p Policy) RequiresPreview(old, next model.Channel) bool {
	return p.RequirePreview && model.IsRiskyChange(old, next)
}
func (p Policy) CanPublish(old, next model.Channel, previewed bool) error {
	if p.RequiresPreview(old, next) && !previewed {
		return errors.New("preview required")
	}
	return p.Check(next)
}
func (p Policy) AllowedTypes() []model.ChannelType {
	return []model.ChannelType{model.Card, model.Wallet, model.Installment, model.Offline}
}
func (p Policy) IsKnown(t model.ChannelType) bool {
	for _, x := range p.AllowedTypes() {
		if x == t {
			return true
		}
	}
	return false
}
func (p Policy) Normalize(c model.Channel) model.Channel {
	if c.FeeBps < 0 {
		c.FeeBps = 0
	}
	if c.LimitCents < 0 {
		c.LimitCents = 0
	}
	return c
}
func (p Policy) Clone() Policy {
	q := New()
	q.RequirePreview = p.RequirePreview
	q.MaxFeeBps = p.MaxFeeBps
	for k, v := range p.Restricted {
		q.Restricted[k] = v
	}
	return q
}
