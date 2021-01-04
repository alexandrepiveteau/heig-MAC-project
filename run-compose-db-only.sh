#!/bin/sh

# Find the right path
cd ./docker/topologies/db || exit 1

# Start Docker Compose
docker-compose down
docker-compose up --build
