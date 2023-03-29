package main

import (
	"fmt"
)

func aFun() {
	var a map[string]string
	// 使用前，需要先make
	a = make(map[string]string, 10)

	a["no1"] = "宋江"
	a["no2"] = "吴用"
	a["no1"] = "武松"
	a["no3"] = "吴用"
	fmt.Println(a)

	for k,v := range a{
		fmt.Printf("\t k=%v, v=%v\n",k,v)
	}
}

func bFun() {
	b := map[int]string{
		1 : "北京",
		2 : "天津",
		3 : "上海",
	}
	
	fmt.Println(b)
	fmt.Println("删除1和4元素")
	delete(b, 1)
	// 不做任何操作
	delete(b, 4)
	fmt.Println(b)
}

func main() {
	aFun()
	bFun()
}