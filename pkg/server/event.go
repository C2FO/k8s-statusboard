package server

import (
	"bytes"
	"fmt"
	"sync"
)

type Eventer interface {
	ToEvent() Event
}

// Event represents a server sent event(SSE)
type Event struct {
	id    string
	event string
	data  []byte
}

// NewEvent initializes an event from its component parts
func NewEvent(id, event string, data []byte) Event {
	return Event{id, event, data}
}

// ToBytes returns the event as a byte slice.
func (e Event) ToBytes() []byte {
	var buf bytes.Buffer
	if e.id != "" {
		buf.WriteString(fmt.Sprintf("id: %s\n", e.id))
	}
	if e.event != "" {
		buf.WriteString(fmt.Sprintf("event: %s\n", e.event))
	}
	buf.WriteString(fmt.Sprintf("data: %s\n", e.data))
	return buf.Bytes()
}

// uniqueEventHolder will be used to store unique events on a per context
// basis.
type uniqueEventHolder struct {
	m      sync.RWMutex
	events map[string]Event
}

func newUniqueEventHolder() *uniqueEventHolder {
	return &uniqueEventHolder{events: make(map[string]Event, 0)}
}

func (h *uniqueEventHolder) AddEvent(e Event) {
	h.m.Lock()
	defer h.m.Unlock()

	h.events[e.event] = e
}

func (h *uniqueEventHolder) Events() []Event {
	h.m.RLock()
	defer h.m.RUnlock()

	events := make([]Event, len(h.events))
	for _, e := range h.events {
		events = append(events, e)
	}
	return events
}
