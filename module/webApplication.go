package module

import (
	"github.com/gin-gonic/gin"
)

func Application(port string) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("0.0.0.0:"+port) // listen and serve on 0.0.0.0:8080
}
