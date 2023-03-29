package main

import (
	"fmt"
)

type Monkey struct {
	Name string
}

// 声明接口
type BridAble interface {
	Flying()
}

type FishAble interface {
	Swimming()
}

func (this *Monkey) climbing() {
	fmt.Println(this.Name, "会爬树..")
}

type LittleMonkey struct {
	Monkey  //继承
}

func (this *LittleMonkey) Flying() {
	fmt.Println(this.Name, "学会了飞翔..")
}

func (this *LittleMonkey) Swimming() {
	fmt.Println(this.Name, "学会了游泳..")
}

func main() {
	monkey := LittleMonkey {
		Monkey {
			Name : "悟空",
		},
	}
	monkey.climbing()
	monkey.Flying()
	monkey.Swimming()
}