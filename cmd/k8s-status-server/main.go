package main

import (
	"github.com/c2fo/k8s-statusboard/pkg/server"
)

func main() {
	s := server.NewStatusServer(8080)
	s.Start()
}
