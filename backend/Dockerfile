FROM golang:alpine AS builder

# for gcc
RUN apk add build-base 

WORKDIR /app/
COPY . .

RUN go build -o server -v .

FROM alpine

RUN apk add --no-cache bash

# RUN apk add curl
COPY --from=builder /app/ /app/
WORKDIR /app/
# add permission to execute for all users and remove all unused docker objects
RUN chmod 777 server && docker system prune -f
# print permissions of server file, was required to test if permission error was from docker container or mac machine of gritlab. It was found to be error coming from bad docker config of mac machine
# CMD ["/bin/ls", "-l"]
CMD ./server

# HEALTHCHECK --interval=3s CMD curl --fail http://localhost:8080/ || exit 1

EXPOSE 8080