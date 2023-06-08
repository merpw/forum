FROM golang:alpine AS builder

ARG SERVICE_NAME=forum

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
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /bin/server backend/${SERVICE_NAME}

FROM alpine

COPY --from=builder /bin/ /bin/

ARG SERVICE_NAME=forum

WORKDIR /app/

VOLUME /app/db

ENV DB_FILE=/app/db/${SERVICE_NAME}.db

CMD server --db=${DB_FILE} --port 8080

EXPOSE 8080