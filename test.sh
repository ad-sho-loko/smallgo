#!/usr/bin/env bash

set -e
gcc -c c/lib.c -o c/lib
files=("basic.go" "var.go" "if.go" "for.go" "array.go" "call.go")
declare -a files

for f in ${files[@]}
do
    echo ${f}
    ./smallgo testdata/${f} > tmp.s
    gcc -g -o tmp tmp.s c/lib
    ./tmp
done

echo All tests success!