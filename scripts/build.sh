#!/bin/bash

IMAGE_NAME="homework-object-storage"

# Function to display messages with a timestamp
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1"
}

log "Starting docker build of $IMAGE_NAME"
docker build -t $IMAGE_NAME ../.
if [ $? -eq 0 ]; then
    log "docker build has finished successfully."
else
    log "Failed to build image."
fi