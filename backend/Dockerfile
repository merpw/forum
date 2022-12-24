FROM golang:alpine AS builder

WORKDIR /app/
COPY . .

RUN go build .

FROM alpine

# RUN apk add curl
# todo audit question with high chance requires that container executable was executed like  "./server", but not an "app/bin/forum"
COPY --from=builder /app/forum /app/bin/
CMD /app/bin/forum

# HEALTHCHECK --interval=3s CMD curl --fail http://localhost:8080/ || exit 1

EXPOSE 8080