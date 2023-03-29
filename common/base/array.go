package main

import (
	"fmt"
	"time"
	"sort"
)

func stocFun() {
	var str string = "hello,world!北京"
	str2 := []rune(str) //把str转成 【】rune

	for i := 0; i < len(str2); i++ {
		fmt.Printf("%c \n", str2[i]);
	}
}

func forFun() {
	var str string = "abc~ok"
	for index, val := range str {
		fmt.Printf("index=%d, val=%c \n", index, val)
	}

	fmt.Printf("倒计时\n")
	j := 0
	for {
		j++
		fmt.Println(j)
		time.Sleep(time.Millisecond * 100)
		if j == 10 {
			break
		}
	}
}

func arrayFun() {
	//var arr [10]int  = [10]int{1, 2, 3, 10, 0, 11, 5, 6, 7, 9}
	var iArr [10]int = [...]int{1, 2, 3, 10, 0, 11, 12, 13, 14, 15}
	fmt.Println("改变前：",iArr)
	iArr[1] = 20
	fmt.Println("改变后：",iArr)

	// 切片
	slice := iArr[2:5]
	fmt.Println("slice 的元素是=",slice)
	fmt.Println("slice 的元素个数", len(slice))
	fmt.Println("slice 容量=", cap(slice))
}

func sortFun() {
	var sort_var []int
	sort_var = append(sort_var, 10)
	sort_var = append(sort_var, 3)
	sort_var = append(sort_var, 7)
	sort_var = append(sort_var, 9)
	fmt.Println("排序前：",sort_var)
	sort.Ints(sort_var)
	fmt.Println("排序后：",sort_var)
}

func main() {
	stocFun()
	forFun()
	arrayFun()
	sortFun()
}