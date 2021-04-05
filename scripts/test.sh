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

    docker network create -d bridge inventory

    export MONGODB_ADMINUSERNAME=root
    export MONGODB_ADMINPASSWORD=example
    export MONGODB_DATABASE="user"

    docker run -d --rm -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME="${MONGODB_ADMINUSERNAME}" -e MONGO_INITDB_ROOT_PASSWORD="${MONGODB_ADMINPASSWORD}" -e MONGO_INITDB_DATABASE="${MONGODB_DATABASE}" --name test-mongo --network=inventory mongo:4.4.4-bionic

    docker run -d --rm -p 6379:6379 --name test-redis --network=inventory redislabs/rejson:1.0.7

    mkdir -p .coverage
    docker build -f Dockerfile.test -t user-service-tests .
    docker run -e MONGODB_ADMINUSERNAME="${MONGODB_ADMINUSERNAME}" -e MONGODB_ADMINPASSWORD="${MONGODB_ADMINPASSWORD}" -e MONGODB_DATABASE="${MONGODB_DATABASE}" --name user-service-tests --network=inventory user-service-tests

    EXIT=$(docker inspect user-service-tests --format='{{.State.ExitCode}}')
    docker cp user-service-tests:/api/coverage.out .coverage/coverage.out
    docker rm user-service-tests
    docker stop test-mongo
    docker stop test-redis
    exit $EXIT
else
    mkdir -p .coverage
    go test ./internal/... -timeout=2m -parallel=4  -covermode=atomic -coverprofile .coverage/coverage.out
fi