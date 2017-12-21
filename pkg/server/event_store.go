package server

import "sync"

var rwmEventStore sync.RWMutex
var eventStore = make(map[string]*uniqueEventHolder, 0)

func addToEventStore(context string, e Event) {
	rwmEventStore.Lock()
	defer rwmEventStore.Unlock()

	if _, ok := eventStore[context]; !ok {
		eventStore[context] = newUniqueEventHolder()
	}
	eventStore[context].AddEvent(e)
}

func getLatestEvents() []Event {
	rwmEventStore.RLock()
	defer rwmEventStore.RUnlock()

	events := make([]Event, 0)
	for _, h := range eventStore {
		events = append(events, h.Events()...)
	}
	return events
}
