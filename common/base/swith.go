// 分支
package main

import (
	"fmt"
)

// if
func ifFun(n int) {
	if n == 0 {
		fmt.Println("------");
	} else if n < 5 {
		fmt.Println("******");
	} else {
		fmt.Println("++++++");
	}
}

// switch
func swFun(n int) {
	switch n {
	case 0:
		fmt.Println("---0---");
	case 6:
		fmt.Println("---6---");
	default:
		fmt.Println("---10---");
	}
}

// 测试
func main() {
	var n int
	fmt.Scanf("输入选择数值：%d",&n)

	fmt.Println("[if测试]")
	ifFun(n)
	fmt.Println("[switch测试]")
	swFun(n)
}