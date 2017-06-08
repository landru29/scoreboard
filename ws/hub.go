package ws

import (
	"encoding/json"
	"errors"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients []*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewHub create a hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run launch the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients = append(h.clients, client)
		case client := <-h.unregister:
			h.removeClient(client)
		case message := <-h.broadcast:
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					h.removeClient(client)
				}
			}
		}
	}
}

// Remove a client from the hub
func (h *Hub) removeClient(c *Client) error {
	index := -1
	for i, client := range h.clients {
		if client == c {
			index = i
		}
	}
	if index < 0 {
		return errors.New("Client not found")
	}
	h.clients = append(h.clients[:index], h.clients[index+1:]...)
	close(c.Send)
	return nil
}

// BroadcastMessage send a message to all clients
func (h *Hub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}

// BroadcastObject send an object as json to all clients
func (h *Hub) BroadcastObject(obj interface{}) (err error) {
	data, err := json.Marshal(obj)
	if err == nil {
		h.BroadcastMessage(data)
	}
	return
}
