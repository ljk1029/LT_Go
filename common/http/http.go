// http 后端处理流程

package main

import (
    "fmt"
    "net/http"
)

// 接收到http请求后处理函数
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// GET 返回信息
		fmt.Fprintln(w, "处理 GET 请求")
        fmt.Println("GET 响应")
	case http.MethodPost:
		// POST
		fmt.Fprintln(w, "处理 POST 请求")
        fmt.Println("POST 响应")
	default:
		// 处理其他请求
		fmt.Fprintln(w, "不支持的请求方法")
        fmt.Println("其他请求 响应")
	}
}

/*
* "/" 这里是请求url后面资源定位  
* 例如客户端url "192.168.10.100:8080/get" 这里可以是"/get"
*/
func main() {
	// 将 handler 函数与根路径 "/" 绑定     
    http.HandleFunc("/", handler)
    fmt.Println("服务端开启") 
   
	// 启动 HTTP 服务器，监听端口 8080
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("Error:", err)
    }
}
