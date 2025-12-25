package websocket

import "time"

type Message struct {
	MatchID int
	Data    []byte
}

// STATE PLAYER TRÊN SERVER
type Player struct {
	ID         uint16
	X, Y       float64
	LastUpdate time.Time
}

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message

	Players map[uint16]*Player //  map lưu state player
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
		Players:    make(map[uint16]*Player), //
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Clients[c] = true

		case c := <-h.Unregister:
			if _, ok := h.Clients[c]; ok {
				delete(h.Clients, c)
				close(c.Send)
			}

		case m := <-h.Broadcast:
			for c := range h.Clients {
				if c.MatchID == m.MatchID {
					select {
					case c.Send <- m.Data:
					default:
						close(c.Send)
						delete(h.Clients, c)
					}
				}
			}
		}
	}
}
