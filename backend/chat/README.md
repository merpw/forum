# backend/chat

Real-time chat service using WebSockets.

## Usage

### Run `go run backend/chat [PARAMS]` to start API server

#### Example: `go run .` to run on default port 8081

### Params:

- `--port` - port to run API server on (default: 8081)
- `--db` - database file path (default: `./chat.db`)

### Environment variables:

- `AUTH_BASE_URL` - optional, default http://localhost:8080 - base url to the auth service

### Testing: `go test backend/chat/... -cover -coverpkg=./backend/chat/...`

### Database migrations

Server supports database migrations. The database will be automatically migrated to the latest version on server start.

If you want to migrate a database manually, you can use `cli`, more info
in [database/migrate](../migrate/README.md).
