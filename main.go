package main

//필요 패키지 가져오기, gorilla, http서버 및 로깅을 위한 패키지
import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
//map에서 웹소켓 연결에 대한 포인터를 정의
var broadcast = make(chan Message)
//클라이언트 측에서 보내온 메시지를 큐잉하는 역할 담당

var upgrader = websocket.Upgrader{}

type Message struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Message string `json:"message"`
}

func main(){
	//채팅 서버의 메인이 되는 변수
	fs := http.FileServer(http.Dir("./front"))
	http.Handle("/", fs)
	
	http.HandleFunc("/ws", handleConnections)
	go handleMessage()
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil{
		log.Fatal("ListenAndServe: ", err)
	}
}

//handleConnections 구현 (연결 핸들러)
//GET 요청을 websocket으로 업그레이드
//받은 요청을 클라이언트로 등록 -> websocket에서 메시지를 기다리고
//받으면 broadcast 채널에 보내는 방식

func handleConnections(responseWriter http.ResponseWriter, request *http.Request){
	//업그레이더 정의. http 연결을 받고 websocket으로 업그레이드 하는 기능을 가진 개체
	var upgrader = websocket.Upgrader{}

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

func handleMessage(){
	for{
		//브로드캐스트 채널에서 다음 메시지를 받는다.
		msg := <- broadcast
		//현재 접속 중인 클라이언트 모두에게 메시지를 보낸다.xw
		for client := range clients{
			err := client.WriteJSON(msg)
			if err != nil{
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
