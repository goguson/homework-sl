#!/bin/bash

# Colima parameters
CPU=1
MEMORY=2
DISK=10

log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1"
}

log "Starting Colima with CPU=$CPU, Memory=${MEMORY}GB, Disk=${DISK}GB..."
colima start --cpu $CPU --memory $MEMORY --disk $DISK
if [ $? -eq 0 ]; then
    log "Colima has started successfully."
else
    log "Failed to start Colima."
fi