FROM golang:1.20-alpine AS builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . ./

RUN go build -v -o clash-exporter .

FROM scratch
WORKDIR /app
COPY --from=builder /app/clash-exporter /app/clash-exporter
ENTRYPOINT ["/app/clash-exporter"]