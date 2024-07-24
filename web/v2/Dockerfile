FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@v0.2.707

COPY . .

RUN templ generate

RUN go build -o nff ./main.go

FROM scratch

COPY --from=builder /app/nff /

CMD [ "./nff" ]
