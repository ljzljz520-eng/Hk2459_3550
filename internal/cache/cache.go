package cache

import (
	"paymentconsole/internal/model"
	"sync"
	"time"
)

type Entry struct {
	Channel model.Channel
	Expires time.Time
}
type Cache struct {
	mu    sync.RWMutex
	items map[string]Entry
	ttl   time.Duration
}

func New(ttl time.Duration) *Cache { return &Cache{items: map[string]Entry{}, ttl: ttl} }
func (c *Cache) Get(id string) (model.Channel, bool) {
	c.mu.RLock()
	e, ok := c.items[id]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.Expires) {
		if ok {
			c.Delete(id)
		}
		return model.Channel{}, false
	}
	return e.Channel, true
}
func (c *Cache) Put(id string, ch model.Channel) {
	c.mu.Lock()
	c.items[id] = Entry{Channel: ch, Expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}
func (c *Cache) Delete(id string)                    { c.mu.Lock(); delete(c.items, id); c.mu.Unlock() }
func (c *Cache) Clear()                              { c.mu.Lock(); c.items = map[string]Entry{}; c.mu.Unlock() }
func (c *Cache) Size() int                           { c.mu.RLock(); defer c.mu.RUnlock(); return len(c.items) }
func (c *Cache) Refresh(id string, ch model.Channel) { c.Put(id, ch) }
func (c *Cache) Expired(id string) bool {
	c.mu.RLock()
	e, ok := c.items[id]
	c.mu.RUnlock()
	return !ok || time.Now().After(e.Expires)
}
