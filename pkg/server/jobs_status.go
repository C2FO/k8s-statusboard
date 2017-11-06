package server

import (
	"encoding/json"

	"k8s.io/client-go/pkg/apis/batch/v1"
)

// JobsStatus is the data structure to hold the data that will be returned for
// all the jobs
type JobsStatus struct {
	Context string   `json:"context"`
	Jobs   []v1.Job `json:"jobs"`
}

// ToEvent creates a (Server Sent) Event for the JobsStatus
func (js JobsStatus) ToEvent() Event {
	b, err := json.Marshal(js)
	if err != nil {
		return Event{}
	}
	return NewEvent("", "job-status", b)
}