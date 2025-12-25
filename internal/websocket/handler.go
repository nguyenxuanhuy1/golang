package websocket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var hub = NewHub()

func init() {
	go hub.Run()
}
func HandleWebSocket(c *gin.Context) {
	matchID, _ := strconv.Atoi(c.Param("match_id"))
	uid, _ := strconv.Atoi(c.Query("user_id"))
	userID := uint16(uid)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(hub, conn, matchID, userID)
	hub.Register <- client

	go client.ReadPump()
	go client.WritePump()

	// gửi event Player Join
	joinPacket := EncodePlayerJoin(userID)
	hub.Broadcast <- Message{
		MatchID: matchID,
		Data:    joinPacket,
	}

	// GỬI VỊ TRÍ BAN ĐẦU CHO CHÍNH NGƯỜI CHƠI (FIX MẤT NHÂN VẬT)
	startMove := EncodeMove(userID, 300, 300)
	client.Send <- startMove // ← FIX CHÍNH

	// rồi mới gửi cho người khác
	hub.Broadcast <- Message{
		MatchID: matchID,
		Data:    startMove,
	}
}
