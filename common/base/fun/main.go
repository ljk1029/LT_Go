package main

import (
	"go/common/base/fun/child"
	"fmt"
)

// func (a bool) test() {
// 	var b fun.MyFamilyAccount
// 	b.loop = a
// 	fmt.Println(b.loop)
// }

func main() {
	fmt.Println("面向对象的方式来完成.....")
	fun.NewMyFamilyAccount().MainMenu()
	//test(true)
}