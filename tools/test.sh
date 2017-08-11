#!/bin/bash

echo "mode: count" > profile.cov

# Standard go tooling behavior is to ignore dirs with leading underscors
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -not -path './vendor/*' -type d);
do
if ls $dir/*.go &> /dev/null; then
    go test -race -covermode=atomic -coverprofile=$dir/profile.tmp $dir
    result=$?
    if [ -f $dir/profile.tmp ]
    then
        cat $dir/profile.tmp | tail -n +2 >> profile.cov
        rm $dir/profile.tmp
    fi
    if [ $result -ne 0 ]; then
        exit $result
    fi
fi
done

go tool cover -func profile.cov
