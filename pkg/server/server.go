package server

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// StatusServer will be responsible for serving our very basic index page,
// its assets, and the broadcasting the server sent events.
type StatusServer struct {
	port   int
	broker *Broker
}

// NewStatusServer creates and initializes a new StatusServer
func NewStatusServer(port int) *StatusServer {
	s := &StatusServer{
		port:   port,
		broker: NewBroker(),
	}
	s.addRoutes()
	s.broker.Start()
	return s
}

func (s *StatusServer) addRoutes() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	http.Handle("/events/", s.broker)
	fs := http.FileServer(http.Dir(wd))
	http.Handle("/", fs)
	// Delegate the events route to our broker which implements ServeHTTP
}

// Start starts the StatusServer
func (s *StatusServer) Start() {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			s.broker.Send(
				NewEvent("", "my-event", []byte("hello")).ToBytes(),
			)
		}
	}()

	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)

	// Loop over our kubernetes clusters and gather information about them every
	// so often so we can send data to our clients.
	// Get some data. Make a new event. and send it.
	// e := NewEvent(...)
	// s.broker.Send(e.ToBytes())
	//
}
