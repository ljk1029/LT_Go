package main

import (
	"fmt"
)

func ifFun(n int) {
	if n == 0 {
		fmt.Println("------");
	} else if n < 5 {
		fmt.Println("******");
	} else {
		fmt.Println("++++++");
	}
}

func swFun(n int) {
	switch n {
	case 0:
		fmt.Println("---2---");
	case 6:
		fmt.Println("---8---");
	default:
		fmt.Println("---10---");
	}
}

func main() {
	var n int
	fmt.Scanf("输入选择数值：%d",&n)

	ifFun(n)
	swFun(n)
}