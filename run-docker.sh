#!/bin/bash

# Set variables
IMAGE_NAME="anne/forumtest"
IMAGE_TAG="1.0"
CONTAINER_NAME="forum_container"
PORT_MAPPING="5000:8080"

# Build the Docker image
echo "Building Docker image..."
docker build -t $IMAGE_NAME:$IMAGE_TAG .

# Check if an old container exists and remove it
if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    echo "Removing old container..."
    docker rm -f $CONTAINER_NAME
fi

# Run the container
echo "Running Docker container..."
docker run -d -p $PORT_MAPPING --name $CONTAINER_NAME $IMAGE_NAME:$IMAGE_TAG

# Show running containers
docker ps
