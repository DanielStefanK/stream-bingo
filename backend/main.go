package main

import (
	"flag"
	"time"

	"github.com/DanielStefanK/stream-bingo/config"
	"github.com/DanielStefanK/stream-bingo/db"
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

	// config path from parameter
	configPath := flag.String("config", "config.yml", "path to config file")

	flag.Parse()

	config.SetConfigPath(*configPath)
	db.Init()

	router := gin.Default()

	// cors
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	auth(router.Group("/auth"))
	router.Run(":" + config.GetConfig().Server.Port)
}

func auth(router *gin.RouterGroup) {
	router.POST("/login", endpoints.Login)
	router.POST("/register", endpoints.Register)
	router.GET("/oauth/:provider", endpoints.OAuthRedirect)
	router.GET("/oauth/:provider/callback", endpoints.OAuthCallback)
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
