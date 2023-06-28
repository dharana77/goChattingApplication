package main

//필요 패키지 가져오기, gorilla, http서버 및 로깅을 위한 패키지
import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

func main(){
	//채팅 서버의 메인이 되는 변수
	var clients = make(map[*websockets.Conn]bool)
	//map에서 웹소켓 연결에 대한 포인터를 정의
	var broadcast = make(chan Message)
	//클라이언트 측에서 보내온 메시지를 큐잉하는 역할 담당
	
	//업그레이더 정의. http 연결을 받고 websocket으로 업그레이드 하는 기능을 가진 개체
	var upgrader = websocket.Upgraders{}
	
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

}

//handleConnections 구현 (연결 핸들러)
//GET 요청을 websocket으로 업그레이드
//받은 요청을 클라이언트로 등록 -> websocket에서 메시지를 기다리고
//받으면 broadcast 채널에 보내는 방식

func handleConnections(responseWriter http.ResponseWriter, request *http.Request){
	
	ws, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil{
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true

	for{
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil{
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}






type Message struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Message string `json:"message"`
}