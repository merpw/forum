#!/bin/bash

echo "emergency.sh begins as docker compose failed due to host misconfiguration"

# go to folder backend
cd backend

# run the command to start the server in the background
go run main.go & # & is used to run the command in the background

# go to folder frontend
cd ../frontend

npm install

# run the command to start the frontend client server
npm run dev

# wait for the terminal disconnect signal
trap 'docker system prune -f; echo "docker pruning started"' SIGHUP

# kill the background process
kill $!
