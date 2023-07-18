#!/bin/bash

BACKUP_DIR=${BACKUP_DIR:-$(pwd)}

docker cp forum-backend-forum:/app/db/database.db "$BACKUP_DIR/database-$(/usr/bin/date "+%Y-%m-%d_%H-%M").db"

docker cp forum-backend-chat:/app/db/chat.db "$BACKUP_DIR/chat-$(/usr/bin/date "+%Y-%m-%d_%H-%M").db"