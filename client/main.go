package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/file", func(c *gin.Context) {
		path := c.Query("path") // shortcut for c.Request.URL.Query().Get("lastname")
		fmt.Println(path);
		c.String(http.StatusOK, "Path received");
	})
	router.Run(":10080")
}