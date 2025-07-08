package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "thaily/proto/asynq"
	"thaily/services/asynq/utils"

	"github.com/hibiken/asynq"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

// WorkflowWorker handles background task processing
type WorkflowWorker struct {
	server  *asynq.Server
	mux     *asynq.ServeMux
	service *AsynqService
}

// NewWorkflowWorker creates a new workflow worker
func NewWorkflowWorker(redisAddr string, redisDB int, service *AsynqService) *WorkflowWorker {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr, DB: redisDB},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				QueueCritical: 6,
				QueueDefault:  3,
				QueueLow:      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	worker := &WorkflowWorker{
		server:  srv,
		mux:     mux,
		service: service,
	}

	// Register handlers
	mux.HandleFunc(TypeWorkflowExecution, worker.HandleWorkflowExecution)

	return worker
}

// Start starts the worker
func (w *WorkflowWorker) Start() error {
	return w.server.Start(w.mux)
}

// Stop stops the worker
func (w *WorkflowWorker) Stop() {
	w.server.Stop()
	w.server.Shutdown()
}

// HandleWorkflowExecution processes workflow execution tasks
func (w *WorkflowWorker) HandleWorkflowExecution(ctx context.Context, t *asynq.Task) error {
	var payload WorkflowTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("Processing workflow execution: %s", payload.ExecutionID)

	// Update task status to RUNNING
	// TODO: Database implementation
	// _, err := w.service.adapter.UpdateOne(ctx, "workflow_executions",
	// 	bson.M{"_id": payload.ExecutionID},
	// 	bson.M{"$set": bson.M{
	// 		"status": "RUNNING",
	// 		"started_at": time.Now(),
	// 	}},
	// )

	// Load workflow from database
	// TODO: Database implementation
	// workflowData, err := w.service.adapter.FindOne(ctx, "workflows",
	// 	bson.M{"_id": payload.WorkflowID},
	// 	bson.M{},
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to load workflow: %w", err)
	// }

	// For now, create a dummy workflow
	workflow := &utils.Workflow{
		ID:   payload.WorkflowID,
		Name: "Test Workflow",
		Steps: []utils.WorkflowStep{
			{
				ID:         "step1",
				Name:       "Test Step",
				ClientName: "auth-service",
				Method:     "Login",
			},
		},
	}

	// Convert input to DynamicInput
	input := &pb.DynamicInput{
		Fields: make(map[string]*structpb.Value),
	}
	for k, v := range payload.Input {
		val, _ := structpb.NewValue(v)
		input.Fields[k] = val
	}

	// Get clients
	clients := make(map[string]*grpc.ClientConn)
	for _, step := range workflow.Steps {
		if client, ok := w.service.GetClient(step.ClientName); ok {
			clients[step.ClientName] = client
		}
	}

	// Execute workflow
	_, err := utils.ExecuteWorkflowSteps(ctx, workflow, input, clients)
	if err != nil {
		// Update status to FAILED
		// TODO: Database implementation
		return fmt.Errorf("workflow execution failed: %w", err)
	}

	// Update status to COMPLETED
	// TODO: Database implementation
	// _, err = w.service.adapter.UpdateOne(ctx, "workflow_executions",
	// 	bson.M{"_id": payload.ExecutionID},
	// 	bson.M{"$set": bson.M{
	// 		"status": "COMPLETED",
	// 		"completed_at": time.Now(),
	// 		"result": result,
	// 	}},
	// )

	log.Printf("Workflow execution completed: %s", payload.ExecutionID)
	return nil
}
