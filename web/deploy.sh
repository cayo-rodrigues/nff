#!/bin/bash

if [ "$#" -lt 4 ]; then
    echo "Usage: $0 <tag_for_web> <tag_for_migrations> <ec2_user_name> <ec2_host_name>"
    exit 1
fi

WEB_IMAGE_TAG=$1
MIGRATIONS_IMAGE_TAG=$2
USER_NAME=$3
HOSTNAME=$4

ssh -o StrictHostKeyChecking=no -i private_key.pem ${USER_NAME}@${HOSTNAME} "
  echo \"version: '3'

  services:
    web:
      container_name: nff-web
      image: ${WEB_IMAGE_TAG}
      environment:
        - PORT=8080
        - DB_URL=postgres://nff:nff@db:5432/nff
        - REDIS_URL=redis://rdb:6379/0
        - SS_API_BASE_URL=https://nff-ss-api-dev.up.railway.app
      ports:
        - 8080:8080
        - 80:8080
      depends_on:
        - db
        - rdb

    db:
      container_name: nff-db
      image: postgres:latest
      ports:
        - 5434:5432
      environment:
        - POSTGRES_USER=nff
        - POSTGRES_PASSWORD=nff
        - POSTGRES_DB=nff
      volumes:
        - postgres_data:/var/lib/postgresql/data

    rdb:
      container_name: nff-rdb
      image: redis:latest
      ports:
        - 6378:6379
      volumes:
        - redis_data:/var/lib/redis/data

    migrations:
      container_name: nff-migrations
      image: ${MIGRATIONS_IMAGE_TAG}
      environment:
        - DB_URL=postgres://nff:nff@db:5432/nff
        - COMMAND=up
      depends_on:
        - db

  volumes:
    postgres_data:
    redis_data:
  \" > docker-compose.yaml

  sudo docker image prune --all --force

  sudo docker pull \"${WEB_IMAGE_TAG}\"
  sudo docker pull \"${MIGRATIONS_IMAGE_TAG}\"

  sudo docker compose run --build --rm migrations

  sudo docker compose down web
  sudo docker image prune --all --force
  sudo docker compose up --build -d web
"
