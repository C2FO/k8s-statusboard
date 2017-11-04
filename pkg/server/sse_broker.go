package server

import (
	"fmt"
	"net/http"
)

// Broker is the message broker that handles broadcasting the messages to each
// client.
type Broker struct {
	// The list of current clients
	clients map[chan []byte]bool
	// The input channel for messages
	in chan []byte
	// Push a channel here to start receiving messages
	connect chan chan []byte
	// Push clients here to disconnect them
	disconnect chan chan []byte
}

// NewBroker creates and initializes a new broker.
func NewBroker() *Broker {
	return &Broker{
		clients:    make(map[chan []byte]bool),
		in:         make(chan []byte),
		connect:    make(chan (chan []byte)),
		disconnect: make(chan (chan []byte)),
	}
}

// Start starts the brokers main loop
func (b *Broker) Start() {
	go func() {
		for {
			select {
			case c := <-b.disconnect:
				delete(b.clients, c)
				close(c)
			case c := <-b.connect:
				b.clients[c] = true
			case msg := <-b.in:
				for c := range b.clients {
					c <- msg
				}
			}
		}
	}()
}

// AddClient adds a client to the broker to start receiving messages
func (b *Broker) AddClient(c chan []byte) {
	b.connect <- c
}

// RemoveClient removes a client from receiving messages.
func (b *Broker) RemoveClient(c chan []byte) {
	b.disconnect <- c
}

// Send sends a message to all of the clients
func (b *Broker) Send(msg []byte) {
	b.in <- msg
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	client := make(chan []byte)
	b.AddClient(client)

	// When the connection is closed, remove the client from the broker
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		b.RemoveClient(client)
	}()

	w.Header().Set("Content-Type", "text/event-stream;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		msg, ok := <-client
		if !ok {
			break // The channel was closed b/c the client disconnected
		}

		fmt.Fprintf(w, "%s\n\n", msg)
		f.Flush()
	}
}
