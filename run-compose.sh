#!/bin/sh

# Find the right path
cd ./docker/topologies/dev || exit 1

# Start Docker Compose
docker-compose down
docker-compose up --build
