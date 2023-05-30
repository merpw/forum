FROM golang:alpine AS builder

# TODO: add universal Dockerfile for all services
# https://github.com/irahardianto/monorepo-microservices/blob/master/Dockerfile

# install gcc, required by sqlite3 driver
RUN apk add build-base

WORKDIR /app/

# Download dependencies first to leverage Docker cache
COPY go.* .
RUN go mod download

COPY . .

# Test with cache
RUN --mount=type=cache,target=/root/.cache/go-build go test ./...

# Build with cache
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /bin/server backend/forum

FROM alpine

COPY --from=builder /bin/ /bin/

WORKDIR /app/

VOLUME /app/db

CMD server --db /app/db/database.db

EXPOSE 8080