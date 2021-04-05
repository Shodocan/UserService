FROM alpine:3

RUN apk update && apk add --no-cache ca-certificates
ADD app /
CMD ["/app"]