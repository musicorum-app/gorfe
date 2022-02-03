FROM golang:latest

USER root
WORKDIR /usr/local/webp
RUN wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-1.0.0.tar.gz \
      && tar -xvzf libwebp-1.0.0.tar.gz \
      && mv libwebp-1.0.0 libwebp && \
      rm libwebp-1.0.0.tar.gz && \
      cd libwebp && \
      ./configure --enable-everything && \
      make && \
      make install && \
      cd .. && \
      rm -rf libwebp

WORKDIR /go/app/bin
ADD . /go
RUN echo "files on /go:" && ls

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
