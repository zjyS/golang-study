package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {

		}
	}(conn)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

// main is the main entry point of the program.
//
// It initializes a router using the gin.Default() function and sets up a GET route at "/ws" with the handler function.
// Then it runs the router on "localhost:8080". If there is an error, it returns.
func main() {
	// 主程序继续执行其他操作
	router := gin.Default()
	router.GET("/ws", handler)
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
