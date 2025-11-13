package websocket

import (
	"encoding/json"
	"log"
	"pongServer/internal/models"
	"sync"
)

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) AddClient(userID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Clients[userID] = client
}

func (h *Hub) RemoveClient(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.Clients, userID)
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			log.Printf("Clinet %s joined the game", client.UserID)
			userConnected := models.WsEvent[models.ConnectedEvent]{
				Type: "user_connected",
				Params: models.ConnectedEvent{
					UserID: client.UserID,
					GameID: client.GameID,
				},
			}
			connectedMsg, err := json.Marshal(userConnected)
			if err != nil {
				log.Println("error marshalling user connected message:", err)
				continue
			}
			client.Send <- connectedMsg
			h.mu.Unlock()
		case client := <-h.Unregister:
			h.mu.Lock()
			_, ok := h.Clients[client.UserID]
			if ok {
				delete(h.Clients, client.UserID)
				log.Printf("Client %s left the game", client.UserID)
				if len(h.Clients) > 0 {
					LeftEventMsg := models.WsEvent[models.LeftEvent]{
						Type: "user_left",
						Params: models.LeftEvent{
							UserID: client.UserID,
							GameID: client.GameID,
						},
					}
					leftMsg, err := json.Marshal(LeftEventMsg)
					if err != nil {
						log.Println("error marshalling user left message:", err)
						continue
					}
					client.Send <- leftMsg
					h.mu.Unlock()
				}
			}

		}
	}
}
