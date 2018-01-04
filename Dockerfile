FROM golang:alpine as build-env
RUN apk add --no-cache git
RUN mkdir -p /go/src/github.com/c2fo/k8s-statusboard
ADD . /go/src/github.com/c2fo/k8s-statusboard
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/c2fo/k8s-statusboard
RUN dep ensure -v
RUN go build -o k8s-statusboard ./cmd/k8s-status-server/main.go

FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/c2fo/k8s-statusboard/k8s-statusboard /app
# Copy the static files to be served by our server
ADD app ./app
# Install the google cloud sdk
RUN apk add --no-cache curl bash python
RUN curl https://sdk.cloud.google.com | bash
EXPOSE 8080
ENTRYPOINT ./k8s-statusboard
