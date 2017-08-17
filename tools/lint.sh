#!/bin/bash

function checkCyclo() {
    cyclo="$(gocyclo $1)"
    complexity="$(echo $cyclo | cut -f 1 -d ' ')"
    if (( complexity > 15)); then
        echo $cyclo
    fi
}

function checkAll() {
    for file in $(find . -name "*.go"); do
        gofmt -s -d $file
        golint $file
        go vet $file
        misspell -error $file
        ineffassign $file
        checkCyclo $file
    done
}

output=$(mktemp)
trap 'rm $output' EXIT

checkAll &> $output

cat $output

if [[ -s $output ]]; then
    exit 1
fi
