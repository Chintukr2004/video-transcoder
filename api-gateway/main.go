package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	
	
	pb "github.com/Chintukr2004/video-transcoder/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	uploadDir := "uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

		
	transcoderHost := os.Getenv("TRANSCODER_HOST")
	if transcoderHost == "" {
		transcoderHost = "localhost:50051"
	}
	conn, err := grpc.NewClient(transcoderHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to transcoder: %v", err)
	}
	defer conn.Close()

	// Create the gRPC client
	transcoderClient := pb.NewTranscoderServiceClient(conn)

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

		// 🚀 NEW: Call the Transcoder Service via gRPC
		// We use a context with a timeout so the gateway doesn't hang forever if the transcoder is busy
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// Determine the absolute path so FFmpeg knows exactly where to look
		absPath, _ := filepath.Abs(savePath)

		// Send the request
		res, err := transcoderClient.StartTranscode(ctx, &pb.TranscodeRequest{
			FilePath:     absPath,
			TargetFormat: "avi", // Let's test converting to AVI first
		})

		if err != nil || !res.Success {
			log.Printf("Transcoding failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "File saved, but transcoding failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     fmt.Sprintf("Success! %s saved and transcoded.", filename),
			"output_path": res.OutputPath,
		})
	})

	fmt.Println("🎬 API Gateway starting on http://localhost:8080...")
	log.Fatal(r.Run(":8080"))
}
