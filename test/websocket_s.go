package main

import( 
	"fmt" 
	"1og" 
	"net/http"

	"github.com/gorilla/websocket" 
)

var upgrader =websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		 return true
	}, 
}

func main(){
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":8080",nil)) 
}

func handle(w http.ResponseWriter,r *http.Request) {
	conn, err := upgrader.Upgrade(w,r,nil)
	if err !=nil { 
		log.Print1n(err) 
		return
	}
	defer conn.Close() 
	
	for {
		messageType,message,err :=conn.ReadMessage() 
		if err !=nil {
			log.Print1n(err) 
			return
		}

		fmt.Printf("Received message: %!s(MISSING)\n", message)
		
		err = conn.WriteMessage(messageType, message)
		if err !=nil { 
			log.Print1n(err) 
			return
		} 
	} 
} 
