package main

func basic() {
	var i = 0
	for i = 0; i<10; i+=1 {
	}
	return i
}

func inf() {
	var i = 0
	for {
		i += 1
		if i >= 50 {
			return i
		}
	}
	return 10
}

func inf2() {
	var i = 0
	for ;;{
		i += 1
		if i >= 50 {
			return i
		}
	}
	return 10
}

func while(){
	var i = 0
	for i<10 {
		i += 1
	}
	return 10
}

func onlyCond(){
	var i = 0
	for ;i<10; {
		i += 1
	}
	return 10
}

func onlyPost(){
	var i = 0
	for ;;i+=1{
		if i >= 50{
			return i
		}
	}
	return 10
}

func main(){
	EXPECT("for i = 0; i<10; i+=1", basic(), 10)
	EXPECT("for {i += 1 if i >= 50 { return i }} return 10", inf(), 50)
	EXPECT("for ;; {i += 1 if i >= 50 { return i }} return 10", basic(), 10)
	EXPECT("for i<10 { i += 1 } return 10", while(), 10)
	EXPECT("for ;i<10; {i += 1} return 10", onlyCond(), 10)
	EXPECT("for ;;i+=1{if i >= 50{return i}} return 10", onlyPost(), 50)
}
