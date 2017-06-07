package ws

import (
	"log"
	"net/http"

	"github.com/satori/go.uuid"
)

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	Log("Connected")
	client := &Client{
		Hub:       hub,
		Connexion: conn,
		Send:      make(chan []byte, 256),
		UUID:      uuid.NewV4().String(),
	}
	client.Hub.register <- client
	go client.writePump()
	go client.readPump()
}
