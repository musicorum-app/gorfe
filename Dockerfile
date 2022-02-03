FROM golang:latest

WORKDIR /go/app/bin
ADD . /go
RUN echo "files on /go:" && ls

USER root

RUN apk add --no-cache \
        freetype-dev \
        libjpeg-turbo-dev \
        libwebp-dev \
        libpng-dev

RUN mkdir -p /go/app/bin
RUN cd /go/src && go get -v && go install -v ./... && go build -o /go/app/bin/goapp
RUN cd /go && ls && cd src && ls
RUN chmod -R 777 /go

RUN mkdir -p /home/musicorum/cache
RUN mkdir -p /home/musicorum/results
RUN chmod -R 777 /home
RUN cd /home/musicorum && ls

EXPOSE 2037

WORKDIR /go/app/bin
ADD . /go/app/bin

ENTRYPOINT /go/app/bin/goapp
