package store

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"paymentconsole/internal/model"
)

func (s *Store) ListAudits() ([]model.AuditEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []model.AuditEvent{}
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(auditBucket).ForEach(func(_, v []byte) error {
			var a model.AuditEvent
			if e := json.Unmarshal(v, &a); e != nil {
				return e
			}
			out = append(out, a)
			return nil
		})
	})
	return out, e
}
func (s *Store) GetWorkflow(id string) (model.Workflow, error) {
	var w model.Workflow
	e := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, workflowBucket, id, &w) })
	return w, e
}
func (s *Store) GetRecord(id string) (model.Record, error) {
	var r model.Record
	e := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, recordBucket, id, &r) })
	return r, e
}
