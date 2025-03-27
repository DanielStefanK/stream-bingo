package main

import (
	"flag"
	"log"
	"time"

	"github.com/DanielStefanK/stream-bingo/auth"
	"github.com/DanielStefanK/stream-bingo/config"
	"github.com/DanielStefanK/stream-bingo/db"
	"github.com/DanielStefanK/stream-bingo/endpoints"
	"github.com/DanielStefanK/stream-bingo/models"
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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.Use(authMiddleware())

	authEndpoints(router.Group("/auth"))
	adminEndpoints(router.Group("/admin"))
	router.Run(":" + config.GetConfig().Server.Port)
}

func authEndpoints(router *gin.RouterGroup) {
	router.POST("/login", endpoints.Login)
	router.POST("/register", endpoints.Register)
	router.GET("/oauth/:provider", endpoints.OAuthRedirect)
	router.GET("/oauth/:provider/callback", endpoints.OAuthCallback)
	router.GET("/me", mustBeAuthenticated(), endpoints.Me)
}

func adminEndpoints(router *gin.RouterGroup) {
	router.Use(mustBeAdmin())
	router.GET("/user/list", endpoints.GetUsers)
	router.DELETE("/user/delete/:userId", endpoints.DeleteUser)
	router.POST("/user/state/:userId", endpoints.DeactiveUser)
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

// middleware for extracting user from jwt and writing it to context
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get token from header
		tokenWithBearer := c.GetHeader("Authorization")
		if tokenWithBearer == "" {
			c.Next()
			return
		}
		token := tokenWithBearer[7:]
		//validate token
		_, err := auth.ValidateJWT(token)
		if err != nil {
			c.Next()
			return
		}

		userId, err := auth.GetUserFromToken(token)
		if err != nil {
			c.Next()
			return
		}

		user := &models.User{}

		db.GetDB().First(user, userId)

		//write user to context
		if user.ID != 0 {
			c.Set("user", user)
		}
		c.Next()
	}
}

// must be authenticated
func mustBeAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exits := c.Get("user")
		if !exits {
			log.Println("User not authenticated", exits)
			c.JSON(401, endpoints.NewErrorResponse(endpoints.ErrAuthorizingUser, "unauthorized", nil))
			c.Abort()
			return
		}
		c.Next()
	}
}

func mustBeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		_, exits := c.Get("user")

		if !exits {
			log.Println("User not authenticated", exits)
			c.JSON(401, endpoints.NewErrorResponse(endpoints.ErrAuthorizingUser, "unauthorized", nil))
			c.Abort()
			return
		}
		user := c.Value("user").(*models.User)
		if user == nil || !user.Admin {
			c.JSON(401, endpoints.NewErrorResponse(endpoints.ErrAuthorizingUser, "unauthorized", nil))
			c.Abort()
			return
		}
		c.Next()
	}
}
