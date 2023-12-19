FROM golang:1.21.5-alpine as builder
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
RUN go build -trimpath -o mqtt2prom .

FROM alpine:3
WORKDIR /app
COPY --from=builder /build/mqtt2prom /bin/mqtt2prom
ENTRYPOINT [ "mqtt2prom" ]
