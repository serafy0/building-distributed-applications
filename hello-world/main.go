package main

import "github.com/gin-gonic/gin"

func IndexHandler(c *gin.Context) {
	name := c.Params.ByName("name")
	c.JSON(200, gin.H{
		"message": "hello " + name,
	})
}

func main() {
	router := gin.Default()

	router.GET("/:name", IndexHandler)
	router.Run()
}
