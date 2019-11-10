package main

func main(){
	var x[10]int
	x[2] = 20
	EXPECT("var x[10] int x[2] = 20 return x[2]", x[2], 20)
}