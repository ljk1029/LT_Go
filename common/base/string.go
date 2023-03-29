package main

import "fmt"
import "strconv"

func inFun() {
	var name string
	var sal float32
	fmt.Println("请输入姓名")
	fmt.Scanln(&name)
	fmt.Println("请输入薪水")
	fmt.Scanln(&sal)
	fmt.Printf("姓名 %v 薪水 %v\n", name, sal)
}

func addrFun() {
	var i1 int = 10
	var ptr *int = &i1
	fmt.Println("i1的地址=", &i1, ptr)
}

func formFun_1() {
	var num1 int = 99
	var num2 float64 = 24.34
	var b bool = true
	var sc byte = 'h'
	var str string

	str = fmt.Sprintf("%d", num1)
	fmt.Printf("str type %T str=%q\n", str, str)

	str = fmt.Sprintf("%f", num2)
	fmt.Printf("str type %T str=%q\n", str, str)

	str = fmt.Sprintf("%t", b)
	fmt.Printf("str type %T str=%q\n", str, str)

	str = fmt.Sprintf("%c", sc)
	fmt.Printf("str type %T str=%q\n", str, str)
}

func formFun_2() {
	var num1 int = 99
	var num2 float64 = 24.34
	var b bool = true
	var str string

	// ways2
	str = strconv.FormatInt(int64(num1), 10)
	fmt.Printf("str type %T str=%q\n", str, str)

	str = strconv.FormatFloat(num2, 'f', 10, 64)
	fmt.Printf("str type %T str=%q\n", str, str)

	str = strconv.FormatBool(b)
	fmt.Printf("str type %T str=%q\n", str, str)

	str = strconv.Itoa(int(num1))
	fmt.Printf("str type %T str=%q\n", str, str)
}

func main() {
	formFun_1()
	formFun_2()
	addrFun()
	addrFun()
}