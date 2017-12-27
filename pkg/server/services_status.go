package server

import (
	"encoding/json"

	"k8s.io/client-go/pkg/api/v1"
)

// ServicesStatus is the data structure to hold the data that will be returned
// for all the services
type ServicesStatus struct {
	Context  string       `json:"context"`
	Services []v1.Service `json:"services"`
}

// ToEvent creates a (Server Sent) Event for the ServicesStatus
func (ss ServicesStatus) ToEvent() Event {
	b, err := json.Marshal(ss)
	if err != nil {
		return Event{}
	}
	return NewEvent("", "services-status", b)
}
