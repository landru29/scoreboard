package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/landru29/scoreboard/src/server/routes/parameters"
	uuid "github.com/satori/go.uuid"
)

const (
	writeWait      = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 1024                // Maximum message size allowed from peer.
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub       *Hub
	Connexion *websocket.Conn // The websocket connection.
	Send      chan []byte     // Buffered channel of outbound messages
	UUID      string          // client UUID
	Role      string          // Client Role
	Offset    int64           // communication offset
}

// Command is the structure to exchange commandes
type Command struct {
	Name      string `json:"name"`
	Data      string `json:"data"`
	Origin    string `json:"origin,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	RequestID string `json:"requestId,omitempty"`
}

// Handshake is a handshake message sent to the client
type Handshake struct {
	Status    string `json:"status"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp,omitempty"`
	RequestID string `json:"requestId,omitempty"`
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// NewHandshake create a handshake
func NewHandshake(status string, data string, requestID string) Handshake {
	return Handshake{
		Timestamp: makeTimestamp(),
		Status:    status,
		Data:      data,
		RequestID: requestID,
	}
}

// NewClient create a newClient on a hub
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	client := &Client{
		Hub:       hub,
		Connexion: conn,
		Send:      make(chan []byte, 1024),
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
		Log("Closing websocket")
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
		switch command.Name {
		case "sync":
			c.Offset = makeTimestamp() - command.Timestamp
			c.Hub.BroadcastObject(NewHandshake("sync", fmt.Sprintf("{\"offset\":%d,\"clientTimestamp\":%d}", c.Offset, command.Timestamp), command.RequestID))
		case "whoami":
			c.SendObject(NewHandshake("whoami", fmt.Sprintf("{\"uuid\":\"%s\"}", c.UUID), command.RequestID))
		case "startjam":
			fmt.Printf("Start Jam\n")
		case "stopJam":
			fmt.Printf("Stop Jam\n")
		case "startTimeout":
			fmt.Printf("Start timeout\n")
		case "stopTimeout":
			fmt.Printf("Stop timeout\n")
		case "updateScore":
			fmt.Printf("Update Score\n")
		case "adjustChronometer":
			fmt.Printf("Adjust chronometere\n")
		case "adjustJam":
			fmt.Printf("Ajust jam\n")
		case "adjustPeriod":
			fmt.Printf("Adjust period\n")
		case "getGameParameters":
			parameter, err := parameters.GetParameter()
			if err != nil {
				c.SendObject(NewHandshake("error", err.Error(), command.RequestID))
				break
			}
			data, err := json.Marshal(parameter)
			if err != nil {
				c.SendObject(NewHandshake("error", err.Error(), command.RequestID))
				break
			}
			fmt.Printf("%s\n", string(data))
			c.Hub.BroadcastObject(NewHandshake("getGameParameters", string(data), command.RequestID))
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
	Log("Will send " + strconv.FormatInt(int64(len(message)), 10) + " bytes")
	Log("Write message " + string(message))
	w.Write(message)

	err = w.Close()
	if err != nil {
		return
	}

	return
}

// SendObject send an object as json
func (c *Client) SendObject(obj interface{}) (err error) {
	bytesToSend, err := json.Marshal(obj)
	if err != nil {
		bytesToSend = []byte("{\"status\":\"error\",\"data\":\"Could not convert response to JSON\"}")
	}
	c.Send <- bytesToSend
	return
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connexion.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// The hub closed the channel.
				c.SendMessage([]byte("{\"status\":\"error\",\"data\":\"Could not transfer the message\"}"))
				return
			}

			c.SendMessage(message)

		case <-ticker.C:
			c.Connexion.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Connexion.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
