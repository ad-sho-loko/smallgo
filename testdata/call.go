package main

func F0() int { var x int x = 10 return x }
func F1(x int) int { return x }
func F3(x, y int) int { return x+y }
func F4(x int, y int, z int) int { return x+y+z }

func main(){
	EXPECT("F0() int { var x int x = 10 return x }", F0(), 10)
	EXPECT("F1(x int) int { return x }", F1(10), 10)
	EXPECT("F2(x, y int) int { return x+y }", F3(10,10), 20)
	EXPECT("F3(x int, y int, z int) int { return x+y+z }", F4(10,10,10), 30)
}
