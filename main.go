package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"gitee.com/Whitroom/imitate-tiktok/sql"
)

func main() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	sql.InitDatabase()

	initRouter(r)

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
