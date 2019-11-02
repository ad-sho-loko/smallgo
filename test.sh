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
try 12 '4*3'
try 11 '2+3*3'
try 2 '6/3'
try 25 '5*(2+3)'
try 6 '3 + +3'
try 0 '3 + -3'
try 1 '1 == 1'
try 0 '1 != 1'

echo OK