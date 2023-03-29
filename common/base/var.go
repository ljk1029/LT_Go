package main

import "fmt"

// 全局变量
var n1 = 100
var nm = "jack"
var(
	n3 = 20
	name2 = "mary"
)

func typeChange() {
	// 类型转换
	var (
		i int32 = 100
		n1 float32 = float32(i)
		n2 int8 = int8(i)
		n3 int64 = int64(1)
	)
	fmt.Printf("i=%v n1=%v n2=%v n3=%v\n", i, n1, n2, n3)
}

func typeFun() {
	//1, 默认值0,再定义i = 1.1 错误，类型不可变
	var i int
	fmt.Println("i=", i)

	//2，自动判定类型
	var num = 0.101
	fmt.Println("num=",num)

	//3, := 左侧不应该是已经声明过的
	name := "tom"
	fmt.Println("name=", name)

	//4, 多个变量申请
	n1, name, n3 := 1000, "tom", 10.233
	fmt.Println("n1=",n1, "name=",name, "n3=",n3)
}

func main() {
	typeFun()
	typeChange()
}