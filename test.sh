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

try 42 'func main() int { return 42 }'
try 5 'func main() int { return 2+3 }'
try 1 'func main() int { return 3-2 }'
try 8 'func main() int { return 2+3+3 }'
try 12 'func main() int { return 4*3 }'
try 11 'func main() int { return 2+3*3 }'
try 2 'func main() int { return 6/3 }'
try 2 'func main() int { return 8%3 }'
try 25 'func main() int {return 5*(2+3) }'
try 6 'func main() int { return 3 + +3 }'
try 0 'func main() int { return 3 + -3 }'
try 1 'func main() int { return 1 == 1 }'
try 0 'func main() int { return 1 != 1 }'
try 1 'func main() int { return 2 > 1 }'
try 1 'func main() int { return 1 < 2 }'
try 0 'func main() int { return 3 < 2 }'
try 1 'func main() int { return 2 >= 2 }'
try 1 'func main() int { return 2 <= 2 }'
try 0 'func main() int { return 3 <= 2 }'
try 2 'func main() int { return 8 >> 2 }'
try 8 'func main() int { return 1 << 3 }'
try 1 'func main() int { return 1+2 > 1*1 }'
try 7 'func main() int { var x = 7 return x }'
try 5 'func main() int { var x = 7 x = 5 return x }'
try 15 'func main() int { var x = 5 var y = 10 return x+y }'
try 8 'func main() int { var x = 2 x <<= 2 return x }'
try 2 'func main() int { var x = 8 x >>= 2 return x }'
try 10 'func main() int { var x = 8 x += 2 return x }'
try 6 'func main() int { var x = 8 x -= 2 return x }'
try 16 'func main() int { var x = 8 x *= 2 return x }'
try 4 'func main() int { var x = 8 x /= 2 return x }'
try 0 'func main() int { var x = 8 x %= 2 return x }'
try 15 'func main() int { var x int x = 15 return x }'
try 15 'func main() int { var x1 int x1 = 15 return x1 }'
try 0 'func main() int { var x int return x }'
try 15 'func main() int { var x int64 x = 15 return x }'

echo OK