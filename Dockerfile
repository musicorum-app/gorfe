FROM golang:alpine AS builder

WORKDIR /go/src
ADD . /go/src

RUN cd /go/src && go install -v ./... && go build -o goapp

EXPOSE 2037
ENTRYPOINT /go/src/goapp
