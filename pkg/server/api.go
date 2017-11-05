package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/c2fo/k8s-statusboard/pkg/k8s"
)

// API implements http.Handler interface for serving the API.
type API struct {
}

// ServeHTTP helps the API implement the http.Handler interface
func (api API) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/contexts":
		api.contexts(rw, r)
	default:
		log.Printf("Unknown route: %s", r.URL.Path)
	}
}

func (api API) contexts(rw http.ResponseWriter, r *http.Request) {
	api.sendJSON(rw, k8s.Contexts())
}

func (api API) sendJSON(rw http.ResponseWriter, v interface{}) {
	rw.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(v)
	if err != nil {
		var buf bytes.Buffer
		rw.WriteHeader(http.StatusInternalServerError)
		buf.WriteString(err.Error())
		rw.Write(buf.Bytes())
	}
}
