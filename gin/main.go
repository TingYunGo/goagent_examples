package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	listen := ":6666"
	router := gin.Default()
	router.GET("/welcome", func(c *gin.Context) {
		fmt.Println("welcom")

		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	fmt.Println("listen ", listen)
	router.Run(listen)
}

