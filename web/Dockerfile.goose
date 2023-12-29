FROM golang:1.21-alpine

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY ./db/migrations/ /migrations

CMD [ "sh", "-c", "goose -dir /migrations -v -s postgres $DB_URL $CMD" ]

