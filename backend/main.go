package main

import (
	"time"

	"github.com/DanielStefanK/stream-bingo/config"
	"github.com/DanielStefanK/stream-bingo/endpoints"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	config.ReloadConfig()
	router := gin.Default()

	auth(router.Group("/auth"))

	router.Run(":8080")
}

func auth(router *gin.RouterGroup) {
	router.POST("/login", endpoints.Login)
	router.POST("/register", endpoints.Register)
	router.GET("/oauth/:provider", endpoints.OAuthCallback)
	router.GET("/oauth/:provider/callback", endpoints.OAuthRedirect)
}

func wsRoutes(router *gin.Engine) {
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
			time.Sleep(time.Second)
		}
	})
}
