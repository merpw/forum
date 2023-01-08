#!/bin/bash

echo "Building development environment..."

echo "Starting backend..."
cd backend && go run . &

echo "Installing frontend dependencies..."
cd frontend && npm install

echo "Starting frontend..."
npm run dev
