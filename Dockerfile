FROM golang:alpine AS builder
WORKDIR /src
COPY ./src /src
RUN cd /src && go build -o goapp
RUN go get

FROM alpine
WORKDIR /app
COPY --from=builder /src/goapp /app

EXPOSE 2037
ENTRYPOINT ./goapp
