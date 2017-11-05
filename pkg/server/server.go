package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/c2fo/k8s-statusboard/pkg/k8s"
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
	http.Handle("/api/", API{})
	http.Handle("/", http.FileServer(http.Dir(wd)))
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
	for _, context := range k8s.Contexts() {
		// Make each request in parallel
		go func(context string) {
			pods, err := k8s.Pods(context)
			if err != nil {
				log.Printf("Error getting pods for %s: %s", context, err)
			}

			ps := PodsStatus{
				Context: context,
				Pods:    pods,
			}
			s.sendEvent(ps.ToEvent())
		}(context)
	}
}

func (s *StatusServer) sendEvent(e Event) {
	s.broker.Send(e.ToBytes())
}
