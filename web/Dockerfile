FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@v0.2.707

COPY . .

RUN templ generate && go mod tidy

RUN go build -o nff ./main.go

RUN apk add --no-cache tzdata

FROM scratch

COPY --from=builder /app/nff /
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

CMD [ "./nff" ]
