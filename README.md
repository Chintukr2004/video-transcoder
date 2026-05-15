# Go Distributed Video Transcoder

An asynchronous, microservices-based video transcoding system built with Go, gRPC, Docker, and FFmpeg.

## 🏗 System Architecture

```mermaid
graph LR
    User([User]) -->|HTTP POST Video| Gateway[API Gateway :8080]
    Gateway -->|Saves Raw File| Vol[(Docker Volume: /uploads)]
    Gateway -->|gRPC: StartTranscode| Worker[Transcoder Node :50051]
    Worker -->|Reads File| Vol
    Worker -->|Executes| FFmpeg[FFmpeg Process]
    FFmpeg -->|Saves AVI| Output[(Local /output)]

🚀 Quick Start
Run the entire distributed system with a single command:

Bash

docker compose up --build

🛠 Tech Stack
Language: Go (Golang)

Networking: gRPC & Protocol Buffers

Web Framework: Gin

Media Processing: FFmpeg (os/exec)

Infrastructure: Docker & Docker-Compose (Multi-stage Alpine builds)