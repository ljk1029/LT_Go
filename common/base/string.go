package main

import "fmt"
import "strconv"

// 输出
func inputFun() {
	var name string
	var sal float32
	fmt.Println("请输入姓名")
	fmt.Scanln(&name)
	fmt.Println("请输入薪水")
	fmt.Scanln(&sal)
	fmt.Printf("姓名 %v 薪水 %v\n", name, sal)
}

// 输出
func outputFun() {
	var i1 int = 10
	var ptr *int = &i1
	fmt.Println("i1的地址=", &i1, ptr)
}

// 字符串格式化，格式成字符串
func formFunA() {
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

// 字符串格式化，格式成字符串
func formFunB() {
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
	fmt.Println("[格式成字符串]")
	formFunA()
	fmt.Println("[格式成字符串]")
	formFunB()
	fmt.Println("[输入测]")
	inputFun()
	fmt.Println("[输出测试]")
	outputFun()
}