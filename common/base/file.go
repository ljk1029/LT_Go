package main

import (
	"fmt"
	"os"
	"bufio"
	"io"
)

var (
	filePath = "/home/li_jk/ljkTest/Go/common/base/test.txt"
	filePath_s = "/home/li_jk/ljkTest/Go/common/base/test.txt"
	filePath_d = "/home/li_jk/ljkTest/Go/common/base/test1.txt"
)


func WriteFun() {
	file, err := os.OpenFile(filePath, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file err =", err)
	}

	defer file.Close()

	// 写
	str_w := "hello,lijk\n"
	writer := bufio.NewWriter(file)
	
	for i:=0; i<5; i++ {
		writer.WriteString(str_w)
	}
	writer.Flush()
}

func ReadFun() {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("open file err =", err)
	}
	// 关闭，防止内存泄露
	defer file.Close()

	const (
		defaultBufSize = 4096
	)
	reader := bufio.NewReader(file)

	// 循环读取文件的内容
	for {
		str,err := reader.ReadString('\n')
		if err == io.EOF { // io.EOF表示文件的末尾
			break
		}
		fmt.Print(str)
	}
	fmt.Print("文件读取结束.....\n")
}

func CopyFile(dstFileName string, srcFileName string)(written int64, err error) {
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		fmt.Printf("open file err=%v\n",err)
	}
	defer srcFile.Close()
	reader := bufio.NewReader(srcFile)

	dstFile, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file err=%v\n",err)
		return
	}
	writer := bufio.NewWriter(dstFile)
	defer dstFile.Close()

	return io.Copy(writer, reader)
}

func main() {
	WriteFun()
	ReadFun()

	// copy
	_, err := CopyFile(filePath_d, filePath_s)
	if err == nil {
		fmt.Printf("拷贝完成\n")
	} else {
		fmt.Printf("拷贝错误 err=%v\n",err)
	}
}