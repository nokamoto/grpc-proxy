#!/bin/bash

set -eu

echo gofmt

if [ -n "$(gofmt -d .)" ]
then
    gofmt -d .
    gofmt -w .
fi

echo golint

golint .

echo go test

go test .
