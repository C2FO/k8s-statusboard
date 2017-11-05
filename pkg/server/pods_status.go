package server

import (
	"encoding/json"

	"k8s.io/client-go/pkg/api/v1"
)

// PodsStatus is the data structure to hold the data that will be returned for
// all the pods
type PodsStatus struct {
	Context string   `json:"context"`
	Pods    []v1.Pod `json:"pods"`
}

// ToEvent creates a (Server Sent) Event for the PodsStatus
func (ps PodsStatus) ToEvent() Event {
	b, err := json.Marshal(ps)
	if err != nil {
		return Event{}
	}
	return NewEvent("", "pod-status", b)
}
