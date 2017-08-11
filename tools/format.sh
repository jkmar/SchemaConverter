#!/bin/bash

for file in $(find . -name "*.go"); do
    gofmt -s -w $file
done
