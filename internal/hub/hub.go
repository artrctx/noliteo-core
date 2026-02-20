package hub

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   uuid.UUID
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) Read() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("websocket error", slog.Any("error", err))
			} else {
				slog.Error("websocket error 2", slog.Any("error", err))
			}
			break
		}
		c.Hub.broadcast <- msg
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		msg, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		// Writer type should be changed depends on use
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			slog.Error("websocket conn next writer errored", slog.Any("error", err))
			return
		}
		if _, err := w.Write(msg); err != nil {
			slog.Error("websocket writer write failed", slog.Any("error", err))
		}

		if err := w.Close(); err != nil {
			slog.Error("websocket writer close failed", slog.Any("error", err))
		}
	}

}

type Hub struct {
	clients    map[*Client]interface{}
	broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func New() *Hub {
	return &Hub{
		clients:    make(map[*Client]interface{}),
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	defer fmt.Println("RUN CLOSED")
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			log.Printf("new client registered %v\n", client.ID)
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Printf("client unregistered %v\n", client.ID)
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.Send <- msg:
				default:
					close(c.Send)
					delete(h.clients, c)
				}
			}
		}
	}
}
