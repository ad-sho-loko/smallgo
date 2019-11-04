#!/usr/bin/env bash

try() {
  expected="$1"
  input="$2"
  # echo "\033[0;31m[test target : $input]\033[0;39m"
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
try 4 'func main() int { return 7 & 4}'
try 9 'func main() int { return 8 | 1}'
try 1 'func main() int { return 3 < 2 || 3 > 2}'
try 0 'func main() int { return 3 < 2 || 3 < 2}'
try 1 'func main() int { return 3 > 2 && 3 > 2}'
try 0 'func main() int { return 2 > 3 && 2 > 3}'
try 1 'func main() int { return 2 >= 2 }'
try 1 'func main() int { return 2 <= 2 }'
try 0 'func main() int { return 3 <= 2 }'
try 2 'func main() int { return 8 >> 2 }'
try 8 'func main() int { return 1 << 3 }'
try 1 'func main() int { return 1+2 > 1*1 }'
try 1 'func main() int { return 1 == 1 > 0 && 2 == 0 <= (2+1) }'
try 11 'func main() int { return 13 ^ 6 }'
try 9 'func main() int { return 13 &^ 6 }'
try 7 'func main() int { var x = 7 return x }'
try 5 'func main() int { var x = 7 x = 5 return x }'
try 15 'func main() int { var x = 5 var y = 10 return x+y }'
try 10 'func main() int { var x = 8 x = x + 2 return x }'
try 10 'func main() int { var x = 8 x += 2 return x }'
try 9 'func main() int { var x = 8 x++ return x }'
try 6 'func main() int { var x = 8 x -= 2 return x }'
try 7 'func main() int { var x = 8 x-- return x }'
try 16 'func main() int { var x = 8 x *= 2 return x }'
try 4 'func main() int { var x = 8 x /= 2 return x }'
try 0 'func main() int { var x = 8 x %= 2 return x }'
try 10 'func main() int { var x = 8 x |= 2 return x }'
try 2 'func main() int { var x = 7 x &= 2 return x }'
try 8 'func main() int { var x = 2 x <<= 2 return x }'
try 2 'func main() int { var x = 8 x >>= 2 return x }'
try 11 'func main() int { var x = 13 x ^= 6 return x }'
try 9 'func main() int { var x = 13 x &^= 6 return x }'
try 15 'func main() int { var x int x = 15 return x }'
try 15 'func main() int { var x1 int x1 = 15 return x1 }'
try 0 'func main() int { var x int return x }'
try 15 'func main() int { var x int64 x = 15 return x }'
try 5 'func main() int { return f() } func f() int { return 5 }'
try 0 'func main() { 5 }'
try 0 'func main() { f() } func f() int { return 5 }'
try 10 'func main() { if 2 > 1 { return 10 } return 5}'
try 5 'func main() { if 2 < 1 { return 10 } return 5}'
try 10 'func main() { var x = 5 if 2 > 1 { x = 10 } return x}'
try 10 'func main() { var x = 5 if 2 > 1 { x = 10 } else { x = 15 } return x}'
try 15 'func main() { var x = 5 if 1 > 2 { x = 10 } else { x = 15 } return x}'
try 5 'func main() { var x = 5 if 1 > 2 { x = 10 } else if 1 > 2 { x = 20 } return x}'
try 20 'func main() { var x = 5 if 1 > 2 { x = 10 } else if 3 > 2 { x = 20 } return x}'
try 30 'func main() { var x = 5 if 1 > 2 { x = 10 } else if 3 > 4 { x = 20 } else { x = 30} return x}'
try 10 'func main() { var i = 0 for i = 0; i<10; i+=1 {} return i }'
try 50 'func main() { var i = 0 for { i += 1 if i >= 50 { return i }} return 10 }'
try 50 'func main() { var i = 0 for ;;{ i += 1 if i >= 50 { return i }} return 10 }'
try 10 'func main() { var i = 0 for i<10 { i += 1 } return 10 }'
try 10 'func main() { var i = 0 for ;i<10; { i += 1 } return 10 }'
try 50 'func main() { var i = 0 for ;;i+=1 { if i >= 50 {return i } } return 10 }'
try 5 'func main() int { return f(5) } func f(n int) int { return 5 }'
try 6 'func main() int { return f(5, 1) } func f(n, m int) int { return n+m }'
try 15 'func main() int { return f(10, 2, 3) } func f(a int, b int, c int) int { return a+b+c }'
try 15 'func main() int { return f(10, 2, 3) } func f(a int, b, c int) int { return a+b+c }'
try 97 "func main() int { var x byte x = 'a' return x}"
try 97 "func main() int { var x byte x = 'a' return F(x) } func F(c byte){ return c }"
try 98 "func main() int { var x byte x = 'a' return F(x) } func F(c byte){ return c+1 }"
try 99 "func main() int { var x byte x = 'a' return F(x, 2) } func F(c byte, i int){ return c+i }"
try 10 'func main() int { var x int x = 10 var y *int y = &x return *y }'
try 10 'func main() int { var x int x = 10 var y *int y = &x var z **int z = &y return **z }'
try 97 "func main() int { var x byte x = 'a' var y *byte y = &x return *y }"
try 1 "func main() int { var x byte x = 'a' var y *byte y = &x return y == y }"

# test errors
# try 10 'func main() { var x = 5 if 2 > 1 { var y = 10 } return y}'

echo OK