FROM golang:1.21 AS builder

WORKDIR /app

RUN apt-get update \
    apt-get install certbot

RUN certbot certonly --standalone -d ec2-18-231-107-141.sa-east-1.compute.amazonaws.com

COPY . .

RUN go build -o nff ./cmd/main.go

FROM scratch

COPY --from=builder /etc/letsencrypt/live/ec2-18-231-107-141.sa-east-1.compute.amazonaws.com/fullchain.pem /
COPY --from=builder /etc/letsencrypt/live/ec2-18-231-107-141.sa-east-1.compute.amazonaws.com/privkey.pem /

COPY --from=builder /app/nff /
COPY --from=builder /app/internal/views /internal/views/
COPY --from=builder /app/internal/static /internal/static/

CMD [ "./nff" ]
