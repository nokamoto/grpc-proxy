#!/bin/bash

set -eu

echo prototool format

prototool format -d examples || prototool format -w examples

echo make examples

for EXAMPLE in {ping,empty-package}
do
    (cd examples/$EXAMPLE; make)
done

echo gofmt

if [ -n "$(gofmt -d .)" ]
then
    gofmt -d .
    gofmt -w .
fi

echo golint

# golint ./... (https://github.com/golang/lint/issues/320)
golint $(go list ./... | grep -v /vendor/)

echo go test

go test ./...
