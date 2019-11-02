#!/usr/bin/env bash

try() {
  expected="$1"
  input="$2"

  echo "[Test : $input]"
  ./smallgo "$input" > tmp.s
  gcc -o tmp tmp.s
  ./tmp
  actual="$?"

  if [[ "$actual" = "$expected" ]]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

try 42 42
try 5 '2+3'
try 1 '3-2'
try 8 '2+3+3'

echo OK