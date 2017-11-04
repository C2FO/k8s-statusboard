package server

import (
	"bytes"
	"fmt"
)

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
