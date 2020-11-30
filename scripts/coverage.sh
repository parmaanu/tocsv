#!/bin/bash
set -e

if [ $# != 1 ]; then
    echo "Usage $0 directory"
    exit
fi

moduledir=$1

echo 'mode: count' > profile.cov

for dir in $(find $moduledir -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d);
do
if ls $dir/*.go &> /dev/null; then
    go test -short -covermode=count -coverprofile=$dir/profile.tmp $dir
    if [ -f $dir/profile.tmp ]
    then
        cat $dir/profile.tmp | tail -n +2 >> profile.cov
        rm $dir/profile.tmp
    fi
fi
done

go tool cover -func profile.cov
go tool cover -html profile.cov -o coverage.html

# mkdir -p cover
# PKG_LIST=`go list ./...`
# for package in ${PKG_LIST}; do
#   go test -covermode=count -coverprofile "cover/${package##*/}.cov" "$package" ;
# done
# tail -q -n +2 cover/*.cov >> cover/coverage.cov
# go tool cover -func=cover/coverage.cov
# go tool cover -html=cover/coverage.cov -o coverage.html
