package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
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

	// client UUID
	UUID string

	// Client Role
	Role string
}

// Command is the structure to exchange commandes
type Command struct {
	Name   string `json:"name"`
	Data   string `json:"data"`
	Origin string `json:"origin,omitempty"`
}

// Handshake is a handshake message sent to the client
type Handshake struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

// NewClient create a newClient on a hub
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	client := &Client{
		Hub:       hub,
		Connexion: conn,
		Send:      make(chan []byte, 256),
		UUID:      uuid.NewV4().String(),
	}
	client.Hub.register <- client

	go client.writePump()
	go client.readPump()

	return client
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Connexion.Close()
	}()
	c.Connexion.SetReadLimit(maxMessageSize)
	c.Connexion.SetReadDeadline(time.Now().Add(pongWait))
	c.Connexion.SetPongHandler(
		func(string) error {
			c.Connexion.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		},
	)

	for {
		_, message, err := c.Connexion.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
				Log("Error " + err.Error())
			}
			break
		}

		fmt.Printf("Raw: %v", message)
		//Log("Raw " + string(message))

		var command Command
		err = json.Unmarshal(message, &command)
		if err != nil {
			c.SendObject(Handshake{
				Status: "error",
				Data:   "This is not a valid command",
			})
		}
		command.Origin = c.UUID

		// Here are the processing of commands

		if c.Hub.BroadcastObject(command) != nil {
			c.SendObject(Handshake{
				Status: "error",
				Data:   "Cannot echo the command to all clients",
			})
		}
	}
}

// SendMessage send a message
func (c *Client) SendMessage(message []byte) (err error) {
	c.Connexion.SetWriteDeadline(time.Now().Add(writeWait))

	w, err := c.Connexion.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	Log("Write message " + strconv.FormatInt(int64(len(message)), 10) + "bytes")
	w.Write(message)

	err = w.Close()
	if err != nil {
		return
	}

	return
}

// SendObject send an object as json
func (c *Client) SendObject(obj interface{}) (err error) {
	data, err := json.Marshal(obj)
	if err == nil {
		return c.SendMessage(data)
	}

	return
}

// writePump pumps messages from the hub to the websocket connection.
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
