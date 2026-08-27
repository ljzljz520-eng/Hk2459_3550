package notify

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	ChannelID string
	Text      string
	At        time.Time
}
type Hub struct {
	mu       sync.Mutex
	messages []Message
}

func New() *Hub { return &Hub{messages: []Message{}} }
func (h *Hub) Publish(id, text string) {
	h.mu.Lock()
	h.messages = append(h.messages, Message{ChannelID: id, Text: text, At: time.Now()})
	h.mu.Unlock()
}
func (h *Hub) List() []Message {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Message, len(h.messages))
	copy(out, h.messages)
	return out
}
func (h *Hub) Latest(id string) (Message, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for i := len(h.messages) - 1; i >= 0; i-- {
		if h.messages[i].ChannelID == id {
			return h.messages[i], true
		}
	}
	return Message{}, false
}
func Format(id string, enabled bool) string {
	state := "disabled"
	if enabled {
		state = "enabled"
	}
	return fmt.Sprintf("channel %s %s", id, state)
}
func (h *Hub) Remove(id string) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	n := 0
	kept := h.messages[:0]
	for _, m := range h.messages {
		if m.ChannelID == id {
			n++
			continue
		}
		kept = append(kept, m)
	}
	h.messages = kept
	return n
}
