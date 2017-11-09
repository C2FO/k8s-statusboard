
.PHONY: build
build:
	go build ./cmd/...

.PHONY: image
image:
	docker build -t k8s-statusboard .

.PHONY: install
install:
	go install ./cmd/...

.PHONY: server
server:
	open http://localhost:8080
	go run ./cmd/k8s-status-server/main.go
