package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"toombs/dataAccess"
)

func main() {
	dataAccess.EstablishConnection()
	defer dataAccess.CloseConnection()

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	albums, err := dataAccess.Albums()
	if err == nil {
		c.IndentedJSON(http.StatusOK, albums)
	}
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum dataAccess.Album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	dataAccess.AddAlbumPtr(&newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	strId := c.Param("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	alb, err := dataAccess.AlbumByID(id)
	if err == nil {
		c.IndentedJSON(http.StatusOK, alb)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
