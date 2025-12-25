package websocket

import (
	"log"
	"time"

	gorilla "github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10

	MapWidth  = 900
	MapHeight = 500
	MaxSpeed  = 300.0 // px/s
)

type Client struct {
	MatchID int
	UserID  uint16
	Hub     *Hub
	Conn    *gorilla.Conn
	Send    chan []byte
}

func NewClient(hub *Hub, conn *gorilla.Conn, matchID int, userID uint16) *Client {
	return &Client{
		MatchID: matchID,
		UserID:  userID,
		Hub:     hub,
		Conn:    conn,
		Send:    make(chan []byte, 256),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		msgType, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType != gorilla.BinaryMessage || len(data) == 0 {
			continue
		}

		event := data[0]

		switch event {

		// Server tự tạo MOVE, nhưng FE có thể gửi MOVE khác?
		// Nếu bạn muốn, có thể bỏ hẳn, nhưng server cần BROADCAST MOVE tạo ra ở handleInput!
		case EventMove:
			// Forward MOVE gói do server tạo sang Hub
			c.Hub.Broadcast <- Message{
				MatchID: c.MatchID,
				Data:    data,
			}

		case EventInput:
			c.handleInput(data)

		case EventShoot:
			c.Hub.Broadcast <- Message{MatchID: c.MatchID, Data: data}

		case EventHit:
			c.Hub.Broadcast <- Message{MatchID: c.MatchID, Data: data}
		}
	}
}

// ================= SERVER-AUTHORITATIVE =================

func (c *Client) handleInput(data []byte) {
	inputMask, err := DecodeInput(data)
	if err != nil {
		log.Println("bad input packet:", err)
		return
	}

	// lấy player state
	p := c.Hub.Players[c.UserID]
	if p == nil {
		p = &Player{
			ID:         c.UserID,
			X:          300,
			Y:          300,
			LastUpdate: time.Now(),
		}
		c.Hub.Players[c.UserID] = p
	}

	// thời gian trôi để tính tốc độ
	now := time.Now()
	dt := now.Sub(p.LastUpdate).Seconds()
	if dt <= 0 {
		dt = 0.016 // fallback 60fps
	}
	p.LastUpdate = now

	// tính di chuyển từ input (UP/DOWN/LEFT/RIGHT)
	vx, vy := 0.0, 0.0

	if inputMask&1 != 0 { // up
		vy -= MaxSpeed * dt
	}
	if inputMask&2 != 0 { // down
		vy += MaxSpeed * dt
	}
	if inputMask&4 != 0 { // left
		vx -= MaxSpeed * dt
	}
	if inputMask&8 != 0 { // right
		vx += MaxSpeed * dt
	}

	// update vị trí
	p.X += vx
	p.Y += vy

	// clamp bản đồ
	if p.X < 0 {
		p.X = 0
	}
	if p.X > MapWidth {
		p.X = MapWidth
	}
	if p.Y < 0 {
		p.Y = 0
	}
	if p.Y > MapHeight {
		p.Y = MapHeight
	}

	// server gửi vị trí thật xuống FE
	safeX := uint16(p.X)
	safeY := uint16(p.Y)

	packet := EncodeMove(p.ID, safeX, safeY)

	c.Hub.Broadcast <- Message{
		MatchID: c.MatchID,
		Data:    packet,
	}
}

// =========================================================

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {

		case msg, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(gorilla.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(gorilla.BinaryMessage, msg)

		case <-ticker.C:
			c.Conn.WriteMessage(gorilla.PingMessage, nil)
		}
	}
}
