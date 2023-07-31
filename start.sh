#!/bin/bash

if [ -z "$1" ]; then
   if [ -z "$TAG" ]; then
      export TAG="latest"
   fi;
else
  export TAG=$1
fi;

export FORUM_BACKEND_SECRET=$(openssl rand -hex 32)

if [ "$TAG" == "latest" ]; then
    echo "Starting latest revision"
    cd production && docker compose up;
    exit
fi;

if [ "$TAG" == "vanilla" ]; then
    echo "Starting vanilla frontend revision"
    docker compose -f docker-compose.vanilla.yml up;
    exit
fi;

if [ "$TAG" == "local" ]; then
    echo "Starting build from local files"
    docker compose up;
    exit
fi;

if [ "$TAG" == "main" ]; then
    echo "Starting main revision"
    docker-compose -f docker-compose.yml -f docker-compose.tag.yml up;
    exit
fi;

if [ -z "$FORUM_TAG" ]; then
    export FORUM_TAG=$TAG
fi;

if [ -z "$CHAT_TAG" ]; then
    export CHAT_TAG=$TAG
fi;

if [ -z "$FRONTEND_TAG" ]; then
    export FRONTEND_TAG=$TAG
fi;

echo "Pulling images forum-backend:$FORUM_TAG, chat-backend:$CHAT_TAG, frontend:$FRONTEND_TAG..."

# Fallback images to main if tag is not found

PULL_RESULT=$(docker compose -f docker-compose.yml -f docker-compose.tag.yml pull 2>&1 | tee /dev/tty \
| grep "Warning" | grep -Eo "backend-forum|backend-chat|frontend" )

while read -r line ; do
    if [ "$line" == "backend-forum" ]; then
      export FORUM_TAG=main
    fi;
    if [ "$line" == "backend-chat" ]; then
      export CHAT_TAG=main
    fi;
    if [ "$line" == "frontend" ]; then
      export FRONTEND_TAG=main
    fi;
done <<< "$PULL_RESULT"

echo "Starting forum-backend:$FORUM_TAG, chat-backend:$CHAT_TAG, frontend:$FRONTEND_TAG..."

docker-compose -f docker-compose.yml -f docker-compose.tag.yml up;