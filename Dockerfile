FROM golang:alpine AS builder

WORKDIR /go
ADD . /go

RUN cd /go/src && go get -v && go install -v ./... && go build -o goapp
RUN chmod -R 777 /go

RUN mkdir -p /home/musicorum/cache
RUN mkdir -p /home/musicorum/results

EXPOSE 2037
ENTRYPOINT /go/src/goapp
