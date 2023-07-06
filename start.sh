#!/bin/bash

TAG=${1:-"latest"}

if [ "$TAG" == "latest" ]; then
    echo "Starting latest revision"
    cd production && FORUM_BACKEND_SECRET=$(openssl rand -hex 32) docker compose up;
    return
fi;

if [ "$TAG" == "local" ]; then
    echo "Starting build from local files"
    FORUM_BACKEND_SECRET=$(openssl rand -hex 32) docker compose up;
    return
fi;

if [ "$TAG" == "main" ]; then
    echo "Starting main revision"
    FORUM_BACKEND_SECRET=$(openssl rand -hex 32) docker-compose -f docker-compose.yml -f docker-compose.tag.yml up;
    return
fi;

# Other revisions, e.g. pr-76. Use `main` tag if the revision doesn't exist
exists () {
    docker pull ghcr.io/merpw/forum/$1 >&2 2>&1
}
tagOrMain () {
    echo "Checking if $1:$TAG exists" >&2
    if exists "$1:$TAG"; then
        echo $TAG
    else
        echo "main"
    fi
}

export FORUM_TAG=$(tagOrMain "backend-forum")
export CHAT_TAG=$(tagOrMain "backend-chat")
export FRONTEND_TAG=$(tagOrMain "frontend")

echo "Starting forum-backend:$FORUM_TAG, chat-backend:$CHAT_TAG, frontend:$FRONTEND_TAG"

TAG=$TAG FORUM_BACKEND_SECRET=$(openssl rand -hex 32) docker-compose -f docker-compose.yml -f docker-compose.tag.yml up;