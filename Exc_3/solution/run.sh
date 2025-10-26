#!/bin/bash
set -e

docker network create orders-net || true

# PostgreSQL 
docker volume create pgdata18 || true

# Stop & remove existing container if exists
docker rm -f orders-db 2>/dev/null || true

docker run -d \
  --name orders-db \
  --network orders-net \
  --env-file debug.env \
  -v pgdata18:/var/lib/postgresql/18/docker \
  postgres:18

# Wait until ready
until docker exec orders-db pg_isready -U docker; do
  sleep 2
done

#  Orderservice build 
docker build -t orderservice .

# Stop & remove existing container if exists
docker rm -f orderservice 2>/dev/null || true

# Orderservice start 
docker run -d \
  --name orderservice \
  --network orders-net \
  --env-file debug.env \
  -p 8080:8080 \
  orderservice
