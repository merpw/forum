#!/bin/bash

BACKUP_DIR=${BACKUP_DIR:-$(pwd)}

docker cp forum-backend:/app/db/database.db "$BACKUP_DIR/database-$(/usr/bin/date "+%Y-%m-%d_%H-%M").db"