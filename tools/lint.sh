#!/bin/bash

for file in $(find . -name "*.go"); do
    golint $file
    go vet $file
    misspell -error $file
    ineffassign $file
done
