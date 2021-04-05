#!/bin/sh

export GO111MODULE=on

# stop old mongo instances

if [ -x "$(command -v docker)" ]
then
    MONGO_DOCKER=$(docker ps | grep test | awk '{print $1}')
    if [  -z "$MONGO_DOCKER" ]
    then
        echo "No Test Running Running"
    else
        docker stop $MONGO_DOCKER
    fi

    export MONGODB_HOST=localhost
    export MONGODB_ADMINUSERNAME=root
    export MONGODB_ADMINPASSWORD=example
    export MONGODB_DATABASE="user"

    docker run -d --rm -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME="${MONGODB_ADMINUSERNAME}" -e MONGO_INITDB_ROOT_PASSWORD="${MONGODB_ADMINPASSWORD}" -e MONGO_INITDB_DATABASE="${MONGODB_DATABASE}" --name test-mongo  mongo:4.4.4-bionic

    export REDISDB_ACTIVE=true
	export REDISDB_HOST=localhost
	export REDISDB_PORT=6379
	export REDISDB_PASSWORD=""
	export REDISDB_DATABASE=1

    docker run -d --rm -p 6379:6379 --name test-redis  redislabs/rejson:1.0.7

    mkdir -p .coverage
    go test ./internal/... -timeout=2m -parallel=4  -covermode=atomic -coverprofile .coverage/coverage.out
    docker stop test-mongo
    docker stop test-redis
else
    mkdir -p .coverage
    go test ./internal/... -timeout=2m -parallel=4  -covermode=atomic -coverprofile .coverage/coverage.out
fi