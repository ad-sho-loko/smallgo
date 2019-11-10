func basicTrue() int {
	if 2 > 1 {
		return 10
	}
	return 5
}

func basicFalse() int {
	if 2 < 1 {
		return 10
	}
	return 5
}

func els() int {
	if 2 < 1 {
		return 10
	} else {
		return 20
	}
	return 5
}

func elseIf() int {
	if 2 < 1 {
		return 10
	} else if 2 > 1{
		return 30
	} else {
		return 20
	}
	return 5
}

func main(){
	EXPECT("if 2 > 1 { return 10 } return 5}", basicTrue(), 10)
	EXPECT("if 2 < 1 { return 10 } return 5}" basicFalse(), 5)
	EXPECT("if 2 < 1 { return 10 } else { return 20 } return 5", els(), 20)
	EXPECT("if 2 < 1 { return 10 } else if 2 > 1 { return 30 } else { return 20 } return 5", elseIf(), 30)
}