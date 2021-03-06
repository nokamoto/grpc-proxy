FROM golang:1.11.0-alpine3.8 AS build

RUN apk update && apk add git

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/nokamoto/grpc-proxy

COPY . .

RUN dep ensure -vendor-only=true

RUN go install .

FROM alpine:3.8

RUN apk update && apk add --no-cache ca-certificates

COPY --from=build /go/bin/grpc-proxy /usr/local/bin/grpc-proxy

ENTRYPOINT [ "grpc-proxy" ]
