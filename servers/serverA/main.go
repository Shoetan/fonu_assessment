package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Shoetan/utils"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

func main() {

	//parent command
	startServer := &cobra.Command{
		Use: "start",
		Short: "A websocket running a TCP server port 8080",
		Run: func(cmd *cobra.Command, args []string) {
			serverStart()
		},
	}

	if err := startServer.Execute(); err != nil {
		log.Fatalf("Command Execution failed %s", err.Error())
	}

	select {} //prevents main from existing


}

func serverStart()  {
		
	log.Println("Starting websocket Server A setup...")


	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

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


	connectServer()


}

func connectServer()  {
	
	conn, err := utils.ConnectToWebSocketServer("localhost:8081")

	if err != nil {
		log.Printf("Could not connect to server B: %v", err.Error())
	}

	log.Printf("Connected Client: %s", conn.RemoteAddr())


	if err !=nil {
		log.Printf("Could not write message to server %s", err.Error())
	}

	exitChan := make(chan bool)

	go func ()  {
		reader := bufio.NewReader(os.Stdin)

		for {

			fmt.Printf("What would you like to do now that you are connected? 😁\n")
			fmt.Printf("1. To send message to server: 1 <add message> 💬\n")
			fmt.Printf(" ctrl + C exits the server 🗑 \n")


			fmt.Println("Enter choice ")
	
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			parts := strings.SplitN(choice, "", 2)// seperate the input from the command line into parts

			var message string 

			if len(parts) >= 2 {
				choice = parts[0] // assign various parts into  respective variables
				message = parts[1]

			}

			switch choice {
			case "1":
				fmt.Println("Message sent ➡ ")
				conn.WriteMessage(websocket.TextMessage, []byte(message))

			case "2":
				fmt.Println("You want to exit the server 🗑")
				exitChan <- true
				return
			}

		}
		
	}()

}