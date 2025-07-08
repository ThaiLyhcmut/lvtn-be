package resolvers

import (
	"context"
	"sync"

	pb "thaily/proto/asynq"

	"thaily/services/adapter"

	"google.golang.org/grpc"
)

type AsynqService struct {
	pb.UnimplementedAsyncQueueServiceServer
	adapter *adapter.MongoDBAdapter

	// gRPC clients map
	clients   map[string]*grpc.ClientConn
	clientsMu sync.RWMutex

	// Redis task client
	redisClient *RedisTaskClient
	
	// Worker for processing tasks
	worker *WorkflowWorker

	// Context for background processing
	ctx    context.Context
	cancel context.CancelFunc
}

type WorkflowTask struct {
	ExecutionID string
	WorkflowID  string
	Input       *pb.DynamicInput
	Context     map[string]string
}

func NewAsynqService(adapter *adapter.MongoDBAdapter, redisAddr string, redisDB int) *AsynqService {
	ctx, cancel := context.WithCancel(context.Background())
	s := &AsynqService{
		adapter:     adapter,
		clients:     make(map[string]*grpc.ClientConn),
		redisClient: NewRedisTaskClient(redisAddr, redisDB),
		ctx:         ctx,
		cancel:      cancel,
	}
	
	// Create worker with service reference
	s.worker = NewWorkflowWorker(redisAddr, redisDB, s)
	
	return s
}

// InitializeClients sets up gRPC client connections
func (s *AsynqService) InitializeClients() error {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	// TODO: Load client configurations from config file or environment
	// For now, using hardcoded example clients
	clientConfigs := map[string]string{
		"auth-service": "localhost:50052",
		"user-service": "localhost:50053",
		// Add more services as needed
	}

	for name, address := range clientConfigs {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return err
		}
		s.clients[name] = conn
	}

	return nil
}

// CloseClients closes all gRPC client connections
func (s *AsynqService) CloseClients() {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	for _, conn := range s.clients {
		conn.Close()
	}
}

// StartTaskProcessor starts the background task processor
func (s *AsynqService) StartTaskProcessor() error {
	return s.worker.Start()
}

// Stop stops the service
func (s *AsynqService) Stop() {
	s.cancel()
	s.worker.Stop()
	s.redisClient.Close()
}


// GetClient returns a gRPC client by name
func (s *AsynqService) GetClient(name string) (*grpc.ClientConn, bool) {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()

	client, ok := s.clients[name]
	return client, ok
}
