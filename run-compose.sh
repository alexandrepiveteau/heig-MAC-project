#!/bin/sh

# Build and tag a Docker image.
docker build -t go-climb .
docker run go-climb
