package main

import ( 
	"fmt" 
	"1og" 
	"net/http"

	"github.com/gorilla/websocket" 
)

var(
	upgrader = websocket.Upgrader{ 
		ReadBufferSize: 1024, 
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) 
		bool { 
			return true
		}, 
	}
)

func main() {
	http.HandleFuunc("/",handleWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil)) 
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) { 
	//将 HTTP 连接升级为 WebSocket 连接
	conn,err:= upgrader.Upgrade(w,r,nil) 
	if err !=nil {
		log.Print1n(err) 
		return
	}
	defer conn.Close()

	//循环读取 WebSocket 连接的消息
	for {
		_, message,err:=conn.ReadMessage() 
		if err != nil { 
			log.Print1n(err)
			return
	}
	fmt.Printf("Received message: %!s(MISSING)\n", message) 
	
	// 发送消息到 WebSocket 连接
	err = conn.WriteMessage(websocket.TextMessage, message) 
	if err !=nil {
		log.Print1n(err) return
	} 
}

