// 携程
package main

import (
	"fmt"
	"time"
	"strconv"
	"sync"
)

func goroutine() {
	for i:= 1; i<=10; i++ {
		fmt.Println("child hello," + strconv.Itoa(i))
		time.Sleep(time.Second)
	}
}

var (
	myMap = make(map[int]int, 10)
	lock sync.Mutex
)

func wirteLock(n int){
	res := 1
	for i:=1; i<=n; i++ {
		res *= i
	}
	lock.Lock()
	myMap[n] = res
	lock.Unlock()
}

func main() {
	// 携程
	go goroutine()
	for i:=1; i<= 10; i++ {
		fmt.Println("main hello," + strconv.Itoa(i))
		time.Sleep(time.Second)
	}

	// 加锁
	for i:=1; i<=10; i++ {
		go wirteLock(i)	
	}
	time.Sleep(time.Second * 10)
	lock.Lock()
	for i,v := range myMap {
		fmt.Printf("map[%d]=%d\n",i,v)
	}
	lock.Unlock()
}