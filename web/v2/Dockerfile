FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o nff ./cmd/web/main.go

FROM scratch

COPY --from=builder /app/nff /

CMD [ "./nff" ]
