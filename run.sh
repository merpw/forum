#!/bin/bash

docker compose -f "docker-compose-vanilla.yml" down

(cd vanilla-frontend && npx tsc)

docker compose -f "docker-compose-vanilla.yml" up --build
