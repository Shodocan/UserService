FROM golang:1.16 AS build
WORKDIR /src

ENV CGO_ENABLED 0
ENV GO111MODULE on

COPY go.mod .
COPY go.sum .

# Enable private libraries installation
RUN go mod download

COPY . .
RUN go build -o "/out/user-service" cmd/main.go 


FROM alpine:latest AS deployment
RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*
WORKDIR /

COPY --from=build /out/user-service /

COPY entrypoint.sh /
RUN chmod +x entrypoint.sh
ENTRYPOINT ./entrypoint.sh user-service