#!/bin/sh
# Start Redis consumer in the background
./redis-consumer &

# Start main server in the foreground
./server
