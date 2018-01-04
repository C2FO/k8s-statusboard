package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/c2fo/k8s-statusboard/pkg/k8s"
)

// StatusServer will be responsible for serving our very basic index page,
// its assets, and the broadcasting the server sent events.
type StatusServer struct {
	broker *Broker
	port   int

	pollingPeriod time.Duration
	eventFuncs    []func(string)
}

// NewStatusServer creates and initializes a new StatusServer
func NewStatusServer(port int) *StatusServer {
	s := &StatusServer{
		port:          port,
		pollingPeriod: 10 * time.Second,
		broker:        NewBroker(),
	}
	// Add event functions that we want to happen every x amount of seconds.
	s.eventFuncs = []func(string){
		s.sendPods,
		s.sendJobs,
		s.sendServices,
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
	http.Handle("/", http.FileServer(http.Dir(filepath.Join(wd, "app"))))
}

// Start starts the StatusServer
func (s *StatusServer) Start() {
	// Go routine to constantly refresh data from all contexts.
	go func() {
		sleepTime := time.Duration(s.pollingPeriod.Nanoseconds() / int64(len(s.eventFuncs)))
		for {
			// For each of the event functions,
			for _, f := range s.eventFuncs {
				for _, context := range k8s.Contexts() {
					go f(context)
				}
				time.Sleep(sleepTime)
			}
		}
	}()

	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *StatusServer) sendPods(context string) {
	pods, err := k8s.Pods(context)
	if err != nil {
		log.Printf("Error getting pods for %s: %s", context, err)
	}

	ps := PodsStatus{
		Context: context,
		Pods:    pods,
	}
	s.updateStoreAndSend(context, ps)
}

func (s *StatusServer) sendJobs(context string) {
	jobs, err := k8s.Jobs(context)
	if err != nil {
		log.Printf("Error getting jobs for %s: %s", context, err)
	}

	js := JobsStatus{
		Context: context,
		Jobs:    jobs,
	}
	s.updateStoreAndSend(context, js)
}

func (s *StatusServer) sendServices(context string) {
	services, err := k8s.Services(context)
	if err != nil {
		log.Printf("Error getting services for %s: %s", context, err)
	}

	ss := ServicesStatus{
		Context:  context,
		Services: services,
	}
	s.updateStoreAndSend(context, ss)
}

func (s *StatusServer) updateStoreAndSend(context string, ei Eventer) {
	e := ei.ToEvent()
	addToEventStore(context, e)
	s.sendEvent(e)
}

func (s *StatusServer) sendEvent(e Event) {
	s.broker.Send(e.ToBytes())
}
