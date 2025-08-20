#!/bin/sh

# Runs PostgreSQL container, which can then be 
# used by the non-dockerized application

. deploy/.env

docker run -it \
    --name postgres \
    -e POSTGRES_DB=${POSTGRES_DB} \
    -e POSTGRES_USER=${POSTGRES_USER} \
    -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
    -v ./deploy/volumes/postgres:/var/lib/postgresql/data \
    -p 127.0.0.1:5432:5432 \
    postgres:17-alpine3.22

docker rm -f postgres >/dev/null