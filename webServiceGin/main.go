package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"toombs/dataAccess"
)

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)

	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	albums, err := dataAccess.Albums()
	if err == nil {
		c.IndentedJSON(http.StatusOK, albums)
	}
}
