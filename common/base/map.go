package main

import (
	"fmt"
)

func mapFunA() {
	var map_a map[string]string
	// 使用前，需要先make
	map_a = make(map[string]string, 10)

	map_a["no1"] = "宋江"
	map_a["no2"] = "吴用"
	map_a["no1"] = "武松"
	map_a["no3"] = "吴用"
	fmt.Println(map_a)

	for k,v := range map_a{
		fmt.Printf("\t k=%v, v=%v\n",k,v)
	}
}

func mapFunB() {
	map_b := map[int]string{
		1 : "北京",
		2 : "天津",
		3 : "上海",
	}
	
	fmt.Println(map_b)
	fmt.Println("删除1和4元素")
	delete(map_b, 1)
	// 不做任何操作
	delete(map_b, 4)
	fmt.Println(map_b)
}

func main() {
	mapFunA()
	mapFunB()
}