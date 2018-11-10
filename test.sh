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

golint . ./yaml

echo go test

go test . ./yaml
