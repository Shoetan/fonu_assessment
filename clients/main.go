package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shoetan/utils"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// connect to the tcp servers
// clients should have the ability to receive messages also not just sent
// This should be the CLI application here
func main()  {

	// set up function call 

	connectToWebsocket := &cobra.Command{
		Use: "connect",
		Short: "Connect to a TCP server",
		Run: func(cmd *cobra.Command, args []string) {
			connectToServer()
		},
	}

	if err := connectToWebsocket.Execute(); err != nil {
		log.Printf("Command Execution failed %s", err.Error())
	}

	select {}


}


func connectToServer()  {

	//connect to the TCP server
	conn, err := utils.ConnectToWebSocketServer("localhost:8081")

	if err != nil {
		log.Fatalf("Failed to connect to server %v", err)
	}


	exitChan := make(chan bool)

	//go routine to receive message from server 
	go func ()  {
			defer conn.Close()
			for {
				_,message, err := conn.ReadMessage() //read message
				if err != nil {
					log.Printf("Error reading from server %v:", err)

					exitChan <- true

					return
				}

				if err == nil && !utils.IsValidMessage(message){
					log.Println("Invalid message type") // check if messagge type is valid or not.
				}

				fmt.Printf("Message recieved from server: %s", string(message))
			}
	}()

	go func ()  {
		reader := bufio.NewReader(os.Stdin)

		for {

			fmt.Printf("What would you like to do now that you are connected? 😁\n")
			fmt.Printf("1. To send broadcast message to server: 1 <add message> 💬\n")
			fmt.Printf("2. To send direct message to server: 2 <add message> 💬\n")
			fmt.Printf(" Ctrl + C exits the server 🗑 \n")


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
				fmt.Println("Sent broadcast message ➡ ")
				conn.WriteMessage(websocket.TextMessage, []byte(message))

			case "2":
				connDetails := conn.LocalAddr()
				fullMessage := fmt.Sprintf("From %s : message %s", connDetails, message)
				fmt.Println("Sent direct message ")
				conn.WriteMessage(websocket.TextMessage, []byte(fullMessage))
			}

		}
		
	}()

}