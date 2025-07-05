package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	pb "thaily/proto/common"
	common "thaily/services/_common/utils"
	"thaily/services/adapter"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		port        = flag.String("port", "50051", "The server port")
		mongoURI    = flag.String("mongo-uri", getEnv("MONGO_URI", "mongodb://thaily:Th%40i2004@localhost:27017"), "MongoDB connection URI")
		mongoDBName = flag.String("mongo-db", getEnv("MONGO_DB", "mongorest"), "MongoDB database name")
	)
	flag.Parse()

	mongoAdapter, err := adapter.NewMongoDBAdapter(*mongoURI, *mongoDBName)
	if err != nil {
		log.Fatalf("Failed to create MongoDB adapter: %v", err)
	}
	defer mongoAdapter.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	service := common.NewCommonService(mongoAdapter)
	pb.RegisterCommonServiceServer(grpcServer, service)

	reflection.Register(grpcServer)

	log.Printf("Common service is running on port %s", *port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
