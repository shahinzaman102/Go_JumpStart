package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader upgrades HTTP requests to WebSocket connections
var upgrader = websocket.Upgrader{}

// Echo is a WebSocket handler that echoes messages back to the client
func Echo(w http.ResponseWriter, r *http.Request) { // Handler for /echo endpoint.
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close() // Ensure connection is closed when function exits

	// Continuously read messages from the client and echo them back
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", msg)

		// Echo the message back to the client
		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}

// WebsocketPage serves the HTML page for the WebSocket frontend
func WebsocketPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/websockets.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}
