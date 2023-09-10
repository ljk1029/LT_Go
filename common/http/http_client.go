// http 客户端

package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

func main() {
    // 服务器地址，包括端口
    serverAddress := "http://localhost:8080"

    // 发送 GET 请求
    response, err := http.Get(serverAddress)
    if err != nil {
        fmt.Println("GET请求失败:", err)
        return
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println("读取响应失败:", err)
        return
    }

    fmt.Println("GET 响应:", string(body))

    // 发送 POST 请求
    postBody := "这是 POST 请求的内容"
    response, err = http.Post(serverAddress, "text/plain", strings.NewReader(postBody))
    if err != nil {
        fmt.Println("POST请求失败:", err)
        return
    }
    defer response.Body.Close()

    body, err = ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println("读取响应失败:", err)
        return
    }

    fmt.Println("POST 响应:", string(body))
}
