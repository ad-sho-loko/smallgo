func main(){
	var x = 10
	EXPECT("var x = 10", x, 10)

	x = 5
	EXPECT("x = 5", x, 5)

	x++
	EXPECT("x++", x, 6)

	x--
	EXPECT("x--", x, 5)

	x+=1
	EXPECT("x+=1", x, 6)

	x-=1
	EXPECT("x-=1", x, 5)

	x*=2
	EXPECT("x*=1", x, 10)

	x/=2
	EXPECT("x/=1", x, 5)

	x%=3
	EXPECT("x%=1", x, 2)

	x|=1
	EXPECT("x|=1", x, 3)

	x&=3
	EXPECT("x&=1", x, 3)

	x<<=1
	EXPECT("x<<=1", x, 6)

	x>>=1
	EXPECT("x>>=1", x, 3)

	x = 13
	x^=6
	EXPECT("x^=1", x, 11)

	x = 13
	x&^=6
	EXPECT("x&^=1", x, 9)

	x = 10
	var y *int
	y = &x

	EXPECT("x = &y return *y", *y, 10)

	var z **int
	z = &y
	EXPECT("z = &y return **z", **z, 10)
}
