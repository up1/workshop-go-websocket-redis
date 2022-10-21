package main

import (
	"demo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	go demo.H.Run()

	router := gin.New()
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	router.GET("/ws/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		demo.ServeWs(c.Writer, c.Request, roomId)
	})
	router.Run("0.0.0.0:8080")
}
