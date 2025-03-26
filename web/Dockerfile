FROM golang:1.24 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y tzdata ca-certificates

RUN go install github.com/a-h/templ/cmd/templ@v0.3.857

COPY . .

RUN templ generate && go mod tidy

RUN go build -o nff

FROM debian:bookworm-slim

WORKDIR /

COPY --from=builder /app/nff /
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/nff"]

