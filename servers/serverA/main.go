package main

import (

	"log"
	"net/http"

	"github.com/Shoetan/utils"
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


}
