package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

func main(){
	var clients = make(map[*websockets.Conn]bool)
	var broadcast = make(chan Message)
}