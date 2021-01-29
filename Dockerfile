FROM golang:alpine AS builder
WORKDIR /src
COPY ./src /src
RUN cd /src && go build -o goapp

FROM alpine
WORKDIR /app
COPY --from=builder /src/goapp /app

EXPOSE 7060
ENTRYPOINT ./goapp
