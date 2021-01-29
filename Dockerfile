FROM golang:alpine AS builder

WORKDIR /go
ADD . /go
RUN echo "files on /go:" && ls

USER root

RUN cd /go/src && go get -v && go install -v ./... && go build -o ../goapp
RUN cd /go && ls && cd src && ls
RUN chmod -R 777 /go

RUN mkdir -p /home/musicorum/cache
RUN mkdir -p /home/musicorum/results
RUN chmod -R 777 /home
RUN cd /home/musicorum && ls

EXPOSE 2037
ENTRYPOINT /go/goapp
