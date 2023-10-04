ARG PORT=8080

FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download && go mod verify

RUN go build -o nff ./cmd/main.go

FROM scratch

COPY --from=builder /app/nff /
COPY --from=builder /app/internal/views /internal/views/
COPY --from=builder /app/internal/static /internal/static/

EXPOSE $PORT

CMD [ "PORT=$PORT ./nff" ]
