FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o nff ./cmd/web/main.go

FROM scratch

COPY --from=builder /app/nff /
COPY --from=builder /app/views /views/
COPY --from=builder /app/static /static/

CMD [ "./nff" ]
