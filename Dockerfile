FROM golang:alpine as builder
COPY . /app
WORKDIR /app
RUN apk update && apk add build-base 
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

COPY /config /config
RUN go install -a -tags "musl,netgo" -ldflags="-s -w" ./cmd/...
