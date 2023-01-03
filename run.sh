#!/bin/bash

# Run docker compose up
docker-compose up &

# Save the PID of the background process
pid=$!

# Wait for the background process to finish
wait $pid


echo "docker-compose up success"
# Wait for the terminal disconnect signal
trap 'docker-compose down --rmi all' SIGHUP
  sleep 5
  docker system prune -f
echo "docker pruning started"
