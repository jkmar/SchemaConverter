#!/bin/bash

for file in $(find . -name "*.go"); do
    gofmt -w $file
done
