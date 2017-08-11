#!/bin/bash

for dir in $(find . -not -path '*/\.*' -type d); do
    echo $dir
    ginkgo -cover $dir
done
