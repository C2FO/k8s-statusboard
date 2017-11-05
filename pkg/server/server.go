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
	broker        *Broker
	port          int
	pollingPeriod time.Duration
}

// NewStatusServer creates and initializes a new StatusServer
func NewStatusServer(port int) *StatusServer {
	s := &StatusServer{
		port:          port,
		pollingPeriod: 15 * time.Second,
		broker:        NewBroker(),
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
	// Delegate the events route to our broker which implements ServeHTTP
	http.Handle("/events/", s.broker)
	fs := http.FileServer(http.Dir(wd))
	http.Handle("/", fs)
}

// Start starts the StatusServer
func (s *StatusServer) Start() {
	go func() {
		for {
			time.Sleep(s.pollingPeriod)
			s.publishK8sStatus()
		}
	}()

	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *StatusServer) publishK8sStatus() {
	// Get k8s clusters.
	// Get Pods stats for each
	// Send events
	s.sendEvent("my-event", []byte("hello"))
}

func (s *StatusServer) sendEvent(event string, data []byte) {
	s.broker.Send(NewEvent("", event, data).ToBytes())
}
