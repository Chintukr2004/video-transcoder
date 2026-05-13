package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	uploadDir := "uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	r := gin.Default()
	r.MaxMultipartMemory = 100 << 20

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("video")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided"})
			return
		}
		filename := filepath.Base(file.Filename)
		savePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Succes! %s saved", filename),
			"path":    savePath,
		})
	})

	fmt.Println("API gateway starting on http://localhost:8080...")
	log.Fatal(r.Run(":8080"))
}
