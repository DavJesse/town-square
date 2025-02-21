#!/bin/bash
echo "Starting a Clean-up ....."
docker container prune -f
docker image prune -af
docker builder prune -af
