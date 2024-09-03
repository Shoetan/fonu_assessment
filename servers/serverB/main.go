package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Shoetan/utils"
	"github.com/gorilla/websocket"
)

func main() {

	log.Println("Starting websocket Server B setup...")

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("New websocket connection attempt on server B")
		utils.HandleConnections(w, r)
	})

	go func ()  {
		log.Println("Websocket server B listening on port :8081")

		err := http.ListenAndServe(":8081", mux)
	
		if err != nil {
			log.Printf("Server encountered an error %v", err)
		}
		log.Println("Websocket server B shut down")
	}()


	time.Sleep(10 * time.Second) // wait for 5 seconds before trying to connect to another websocket



	conn, err := utils.ConnectToWebSocketServer("localhost:8080")

	if err != nil {
		log.Printf("Could not connect to server A: %v", err.Error())
	}

	log.Printf("Connected Client: %s", conn.RemoteAddr().Network())

	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello to another socket"))

	if err !=nil {
		log.Printf("Could not write message to server %s", err.Error())
	}

	select{}


}