package store

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"paymentconsole/internal/model"
	"sync"
	"time"
)

var channelsBucket = []byte("channels")
var auditBucket = []byte("audit")
var workflowBucket = []byte("workflows")
var recordBucket = []byte("records")

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	db, e := bbolt.Open(path, 0600, nil)
	if e != nil {
		return nil, e
	}
	s := &Store{db: db}
	e = db.Update(func(tx *bbolt.Tx) error {
		for _, b := range [][]byte{channelsBucket, auditBucket, workflowBucket, recordBucket} {
			if _, x := tx.CreateBucketIfNotExists(b); x != nil {
				return x
			}
		}
		return nil
	})
	if e != nil {
		db.Close()
		return nil, e
	}
	return s, nil
}
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
func putJSON(tx *bbolt.Tx, b []byte, key string, v any) error {
	raw, e := json.Marshal(v)
	if e != nil {
		return e
	}
	return tx.Bucket(b).Put([]byte(key), raw)
}
func getJSON(tx *bbolt.Tx, b []byte, key string, v any) error {
	raw := tx.Bucket(b).Get([]byte(key))
	if raw == nil {
		return bbolt.ErrBucketNotFound
	}
	return json.Unmarshal(raw, v)
}
func (s *Store) SaveChannel(c model.Channel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, channelsBucket, c.ID, c) })
}
func (s *Store) GetChannel(id string) (model.Channel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var c model.Channel
	e := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, channelsBucket, id, &c) })
	return c, e
}
func (s *Store) ListChannels() ([]model.Channel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []model.Channel{}
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(channelsBucket).ForEach(func(_, v []byte) error {
			var c model.Channel
			if e := json.Unmarshal(v, &c); e != nil {
				return e
			}
			out = append(out, c)
			return nil
		})
	})
	return out, e
}
func (s *Store) SaveAudit(a model.AuditEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, auditBucket, a.ID, a) })
}
func (s *Store) SaveWorkflow(w model.Workflow) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, workflowBucket, w.ID, w) })
}
func (s *Store) SaveRecord(r model.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, recordBucket, r.ID, r) })
}
func (s *Store) Health() error { return s.db.View(func(*bbolt.Tx) error { return nil }) }
func (s *Store) Count(bucket []byte) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	s.db.View(func(tx *bbolt.Tx) error { n = tx.Bucket(bucket).Stats().KeyN; return nil })
	return n
}
func ID(prefix string) string { return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano()) }
