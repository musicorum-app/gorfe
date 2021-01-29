FROM golang:alpine AS builder
WORKDIR /build
COPY ./ /build
RUN mkdir /out
RUN cd /build/src && go get && build -o /out/goapp

FROM alpine
WORKDIR /app/gorfe
COPY --from=builder /out /app/gorfe
COPY --from=builder /build/src/vendor /app/gorfe/vendor

EXPOSE 2037
ENTRYPOINT ./goapp
