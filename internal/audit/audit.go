package audit

import (
	"fmt"
	"paymentconsole/internal/model"
	"paymentconsole/internal/store"
)

type Logger struct{ s *store.Store }

func New(s *store.Store) *Logger { return &Logger{s: s} }
func (l *Logger) Record(action, channel, actor string, before, after bool) error {
	return l.s.SaveAudit(model.AuditEvent{ID: store.ID("audit"), Action: action, ChannelID: channel, Actor: actor, Before: before, After: after})
}
func (l *Logger) Describe(a model.AuditEvent) string {
	return fmt.Sprintf("%s %s %t->%t", a.Actor, a.ChannelID, a.Before, a.After)
}
