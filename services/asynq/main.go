package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	pb "thaily/proto/asynq"
	"thaily/services/adapter"
	"thaily/services/asynq/resolvers"

	"google.golang.org/grpc"
)

var (
	port     = flag.String("port", "50055", "The server port")
	mongoURI = flag.String("mongo-uri", getEnv("MONGO_URI", "mongodb://localhost:27017"), "MongoDB URI")
	dbName   = flag.String("db-name", getEnv("DB_NAME", "asynq"), "Database name")
	redisAddr = flag.String("redis-addr", getEnv("REDIS_ADDR", "localhost:6379"), "Redis address")
	redisDB   = flag.Int("redis-db", 0, "Redis database number")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	flag.Parse()

	// Initialize MongoDB adapter
	mongoAdapter, err := adapter.NewMongoDBAdapter(*mongoURI, *dbName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoAdapter.Close()

	// Create asynq service with Redis
	asynqService := resolvers.NewAsynqService(mongoAdapter, *redisAddr, *redisDB)

	// Initialize gRPC clients from config
	if err := asynqService.InitializeClients(); err != nil {
		log.Fatalf("Failed to initialize clients: %v", err)
	}
	defer asynqService.CloseClients()

	// Start async task processor
	if err := asynqService.StartTaskProcessor(); err != nil {
		log.Fatalf("Failed to start task processor: %v", err)
	}
	defer asynqService.Stop()

	// Create gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAsyncQueueServiceServer(grpcServer, asynqService)

	log.Printf("Asynq service listening on port %s", *port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
