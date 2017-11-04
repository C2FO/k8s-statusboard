package server

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
