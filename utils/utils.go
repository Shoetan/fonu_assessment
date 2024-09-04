package utils

import (
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager holds the WebSocket connections
type ConnectionManager struct {
	connections map[string]*websocket.Conn
	mu          sync.Mutex
}



func LogError(err error, message string)  {

	if err != nil {
		log.Printf("%s : %v", message, err.Error())
	}

}

// websocket upgrader upgrades the standard http connection to websocket
var upgrader =  websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,

	// check origin of connection allows all connections
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // upgrades the http server to a websocket

	LogError(err, "Could not upgrade http server to websockets")

	log.Printf("Client connected : %s", conn.RemoteAddr().String())

	for {

		messageType, message, err := conn.ReadMessage()
	
		if err != nil {
			log.Printf("Could not read message from connection %s", err.Error())
			return
		}
	
		log.Printf("MessageType: %v, Received message from client: %s, message: %s", messageType,conn.RemoteAddr() ,string(message))

		// send message back to sender
		err = conn.WriteMessage(messageType, message)

		if err != nil {
			log.Println("Could not send the message back")
		}


	}

}

func ConnectToWebSocketServer(serverURL string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: serverURL, Path: "/ws"} //parse url 

	log.Printf("Connecting to %s", u.String())

	conn,_,err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		return nil, err
	}

	return conn, nil

}

func NewConnectionManager() *ConnectionManager{
	return &ConnectionManager{
		connections: make(map[string]*websocket.Conn),
	}
}

//Adding a new websocket connetion to map 

func(cm *ConnectionManager) AddConnection(id string, conn *websocket.Conn) {
	cm.mu.Lock()
	cm.connections[id] =  conn
	cm.mu.Unlock()

}

//Removing connections from map
func (cm * ConnectionManager) RemoveConnections(id string)  {
	cm.mu.Lock()
	delete(cm.connections, id)
	cm.mu.Unlock()
}
