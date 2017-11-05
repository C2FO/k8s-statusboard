package main

import (
	"github.com/c2fo/k8s-statusboard/pkg/k8s"
	"github.com/c2fo/k8s-statusboard/pkg/server"
)

func main() {
	k8s.LoadConfig()
	s := server.NewStatusServer(8080)
	s.Start()
}
