package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Shoetan/utils"
	"github.com/gorilla/websocket"
)

func main() {

	log.Println("Starting websocket Server A setup...")

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("New websocket connection attempt on server A")
		utils.HandleConnections(w, r)
	})

	go func ()  {
		log.Println("Websocket server A listening on port :8080")

		err := http.ListenAndServe(":8080", mux)
	
		if err != nil {
			log.Printf("Server encountered an error %v", err)
		}
		log.Println("Websocket server A shut down")
	}()

	time.Sleep(10 * time.Second) // wait for 5 seconds before trying to connect to another websocket


	conn, err := utils.ConnectToWebSocketServer("localhost:8081")

	if err != nil {
		log.Printf("Could not connect to server A: %v", err.Error())
	}

	log.Printf("Connected Client: %s", conn.RemoteAddr())

	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello to socket go routine"))

	if err !=nil {
		log.Printf("Could not write message to server %s", err.Error())
	}

	select{}

}