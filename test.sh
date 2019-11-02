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

try 42 'return 42'
try 5 'return 2+3'
try 1 'return 3-2'
try 8 'return 2+3+3'
try 12 'return 4*3'
try 11 'return 2+3*3'
try 2 'return 6/3'
try 2 'return 8%3'
try 25 'return 5*(2+3)'
try 6 'return 3 + +3'
try 0 'return 3 + -3'
try 1 'return 1 == 1'
try 0 'return 1 != 1'
try 1 'return 2 > 1'
try 1 'return 1 < 2'
try 0 'return 3 < 2'
try 1 'return 2 >= 2'
try 1 'return 2 <= 2'
try 0 'return 3 <= 2'
try 2 'return 8 >> 2'
try 8 'return 1 << 3'
try 1 'return 1+2 > 1*1'
try 7 'var x = 7 return x'
try 15 'var x = 5 var y = 10 return x+y'
# try 15 'var x int x = 15 return x'

echo OK