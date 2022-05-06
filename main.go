package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	initRouter(r)

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
