package ws

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Connexion *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	UUID string
}

// Command is the structure to exchange commandes
type Command struct {
	Name   string `json:"name"`
	Data   string `json:"data"`
	Origin string `json:"origin"`
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Connexion.Close()
	}()
	c.Connexion.SetReadLimit(maxMessageSize)
	c.Connexion.SetReadDeadline(time.Now().Add(pongWait))
	c.Connexion.SetPongHandler(func(string) error { c.Connexion.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Connexion.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
				Log("Error " + err.Error())
			}
			break
		}
		Log("Read message " + strconv.FormatInt(int64(len(message)), 10) + "bytes")
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.Hub.broadcast <- message
	}
}

// SendMessage send a message
func (c *Client) SendMessage(message []byte) {
	c.Connexion.SetWriteDeadline(time.Now().Add(writeWait))

	w, err := c.Connexion.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	Log("Write message " + strconv.FormatInt(int64(len(message)), 10) + "bytes")
	w.Write(message)

	if err := w.Close(); err != nil {
		return
	}

}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	command := Command{
		Name:   "SetUuid",
		Data:   c.UUID,
		Origin: "SERVER",
	}
	data, err := json.Marshal(command)
	if err != nil {
		c.SendMessage([]byte("{\"status\":\"fatal\"}"))
	} else {
		c.SendMessage(data)
	}

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connexion.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Connexion.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Connexion.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Connexion.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			Log("Write message " + strconv.FormatInt(int64(len(message)), 10) + "bytes")
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Connexion.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Connexion.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
