package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	pb "github.com/Chintukr2004/video-transcoder/proto"
	"google.golang.org/grpc"
)

// server is used to implement the TranscoderService defined in proto
type server struct {
	pb.UnimplementedTranscoderServiceServer
}

// StartTranscode is the function that gets called by the API Gateway via gRPC
func (s *server) StartTranscode(ctx context.Context, req *pb.TranscodeRequest) (*pb.TranscodeResponse, error) {
	log.Printf(" Received request to transcode: %s to format: %s\n", req.FilePath, req.TargetFormat)

	// Create an output directory if it doesn't exist
	outDir := "output"
	os.MkdirAll(outDir, os.ModePerm)

	// Determine the output file name (e.g., test_video.mp4 -> test_video_converted.avi)
	filename := filepath.Base(req.FilePath)
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[0 : len(filename)-len(ext)]
	outputPath := filepath.Join(outDir, fmt.Sprintf("%s_converted.%s", nameWithoutExt, req.TargetFormat))


	go func(inputFile, outputFile string){
		log.Printf("Starting background transcode for %s...", inputFile)
		cmd := exec.Command("ffmpeg", "-y", "-i", inputFile, "-preset","ultrafast", outputFile)
		output, err := cmd.CombinedOutput()

		if err != nil{
			log.Printf("Background FFmpeg Error for %s: %v\nOutput: %s\n", inputFile,err, string(output))
			return
		}
		log.Printf("Background transcode complete! Saved to: %s\n", outputFile)
	}(req.FilePath, outputPath)
	return &pb.TranscodeResponse{
		Success: true,
		Message: "video added to proccesing queue",
		OutputPath: "Proccessing in background",
	}, nil 
}
	

func main() {
	//  Listen on a specific port for internal gRPC traffic
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	
	grpcServer := grpc.NewServer()

	// Register our transcoder service with the gRPC server
	pb.RegisterTranscoderServiceServer(grpcServer, &server{})

	fmt.Println("Transcoder Service listening on port :50051...")
	

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}