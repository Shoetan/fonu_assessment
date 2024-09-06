package utils

import (
	// "fmt"
	"log"
	// "net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"regexp"
	"github.com/gorilla/websocket"
)

// ConnectionManager holds the WebSocket connections
type ConnectionManager struct {
	connections map[string]*websocket.Conn
	mu          sync.RWMutex
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

var cm = NewConnectionManager() //Global connection manager

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // upgrades the http server to a websocket

	var ipAddress string 

	LogError(err, "Could not upgrade http server to websockets")

	cm.AddConnection(conn.RemoteAddr().String(), conn) // add connection to connection pool

	for {

		_, message, err := conn.ReadMessage()
	
		if err != nil {
			log.Printf("Could not read message from connection %s", err.Error())
			return
		}

		messageString := string(message)
	
		log.Printf("Received message from client: %s, message: %s",conn.RemoteAddr() ,messageString)
		log.Println(cm.connections)

		if strings.Contains(messageString, "From") {
			//Extract the senderID and the actual message
			parts := strings.SplitN(messageString, " : ", 2)

			if len(parts) == 2 {
				from := parts[0]
				re := regexp.MustCompile(`From\s*(\[[^\]]+\]:\d+)`)
					// Find the first match
				match := re.FindStringSubmatch(from)

				if len(match) > 1 {
					ipAddress = match[1]
				} 

				msg := parts[1]
				cm.SendDirectMessage(ipAddress, []byte(msg))
			} 
		}
		cm.BroadcastMessage([]byte(message))


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
	log.Printf("Client: %s", conn.RemoteAddr().String())
	cm.mu.Unlock()

}

//Removing connections from map
func (cm * ConnectionManager) RemoveConnections(id string)  {
	cm.mu.Lock()
	delete(cm.connections, id)
	cm.mu.Unlock()
}

func (cm *ConnectionManager) BroadcastMessage(message []byte) {
	cm.mu.RLock()
	conns := make([]*websocket.Conn, 0, len(cm.connections)) // fetch connections first 

	for _, conn := range cm.connections{
		conns = append(conns, conn)
	}
	cm.mu.RUnlock() //clean up later

	for _, conn := range conns {
		err := conn.WriteMessage(websocket.TextMessage, message)

		if err != nil {
			log.Printf("Error broadcasting message: %v", err)
		}
	}
}

func (cm * ConnectionManager) SendDirectMessage(id string, message []byte)  {
	cm.mu.RLock()

	// DirectMessage(message, conn) // collect the connection as string and try to find it"s match in the client pool

	//fetch connections first too 
	connections := make(map[string]*websocket.Conn)
	for addr, conn := range cm.connections {
		connections[addr] = conn
	}

	cm.mu.RUnlock()

	// Retrieve the connection associated with the given ID
	storedConn, ok := cm.connections[id]
	if !ok {
		log.Printf("connection not found for address %s", id)
	}

	err := storedConn.WriteMessage(websocket.TextMessage, message)

	if err != nil {
		log.Printf("Could not send direct message %s", err.Error())
	}

}

func DirectMessage(message []byte, conn *websocket.Conn)  {
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Error sending direct message to server %s: Error: %s", conn.LocalAddr(), err )
	}
}
