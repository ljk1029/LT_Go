func getVal(num1 int, num2 int)(int, int)
{
	sum := num1 + num2
	sub := num2 - num1
	return sum, sub
}

func mian()
{
	sum, sub := getVal(30, 15)
	sum2, -  := getVal(20, 18)  // 只取第一个值
}