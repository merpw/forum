#!/bin/bash

(cd vanilla-frontend && npx tsc -w) & (cd backend/forum && go run .) & (cd backend/chat && go run .)
