package router

import (
	"traingolang/internal/api/handler"
	// "traingolang/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// User
	r.POST("/user/register", handler.Register)
	r.POST("/user/login", handler.Login)
	// Match API
	// r.POST("/match/create", handler.CreateMatch)
	// r.POST("/match/join", handler.JoinMatch)

	// // WebSocket (phải có match_id)
	// r.GET("/ws/:match_id", websocket.HandleWebSocket)

	return r
}
