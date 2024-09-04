package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	// "github.com/Shoetan/utils"
)



func TestConnectToWebSocketServer(t *testing.T) {
	// Start a test WebSocket server
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		_, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade websocket: %v", err)
		}
	}))
	defer s.Close()

	// Extract the server URL
	serverURL := s.Listener.Addr().String()

	// Test the connection
	conn, err := ConnectToWebSocketServer(serverURL)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	// Test with an invalid URL
	_, err = ConnectToWebSocketServer("invalid_url")
	assert.Error(t, err)

	// Close the connection
	if conn != nil {
		conn.Close()
	}
}

func TestHandleConnections(t *testing.T) {
	// Start a test HTTP server with the HandleConnections handler
	s := httptest.NewServer(http.HandlerFunc(HandleConnections))
	defer s.Close()

	// Convert the test server's URL to WebSocket format
	wsURL := "ws" + s.URL[len("http"):]

	// Dial the server using WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err, "Failed to connect to WebSocket server")
	defer ws.Close()

	// Test sending a message to the server
	testMessage := "Hello, WebSocket!"
	err = ws.WriteMessage(websocket.TextMessage, []byte(testMessage))
	assert.NoError(t, err, "Failed to send message to WebSocket server")

	// Test receiving the echoed message from the server
	_, message, err := ws.ReadMessage()
	assert.NoError(t, err, "Failed to read message from WebSocket server")
	assert.Equal(t, testMessage, string(message), "Received message does not match sent message")
}

