package model

import "time"

type ChannelType string

const (
	Card        ChannelType = "card"
	Wallet      ChannelType = "wallet"
	Installment ChannelType = "installment"
	Offline     ChannelType = "offline"
)

type Channel struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Type       ChannelType `json:"type"`
	Enabled    bool        `json:"enabled"`
	LimitCents int64       `json:"limit_cents"`
	FeeBps     int         `json:"fee_bps"`
	Notice     string      `json:"notice"`
	Version    int64       `json:"version"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
type PaymentRequest struct {
	ID          string
	ChannelID   string
	AmountCents int64
	CreatedAt   time.Time
}
type PaymentResult struct {
	ID        string
	ChannelID string
	Accepted  bool
	Reason    string
	FeeCents  int64
}
type AuditEvent struct {
	ID        string
	Action    string
	ChannelID string
	Before    bool
	After     bool
	Actor     string
	CreatedAt time.Time
}
type Workflow struct {
	ID        string
	Name      string
	Steps     []string
	Status    string
	CreatedAt time.Time
}
type Record struct {
	ID        string
	Kind      string
	Payload   []byte
	CreatedAt time.Time
}
type Attachment struct {
	ID         string
	WorkflowID string
	Name       string
	Data       []byte
}

func NewChannel(id string, t ChannelType) Channel {
	return Channel{ID: id, Name: string(t), Type: t, Enabled: true, UpdatedAt: time.Now()}
}
func (c Channel) CanAccept(amount int64) bool {
	if !c.Enabled {
		return false
	}
	if amount <= 0 {
		return false
	}
	if c.LimitCents > 0 && amount > c.LimitCents {
		return false
	}
	return true
}
func (c Channel) Fee(amount int64) int64 {
	if amount < 0 {
		return 0
	}
	return amount * int64(c.FeeBps) / 10000
}
func (c Channel) Clone() Channel { return c }
