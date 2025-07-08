package resolvers

import (
	"context"
	"fmt"

	pb "thaily/proto/asynq"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ExecuteWorkflow implements the ExecuteWorkflow RPC
func (s *AsynqService) ExecuteWorkflow(ctx context.Context, req *pb.ExecuteWorkflowRequest) (*pb.ExecuteWorkflowResponse, error) {
	// Validate request
	if req.WorkflowId == "" {
		return nil, status.Error(codes.InvalidArgument, "workflow_id is required")
	}

	// Check if workflow exists in database
	// TODO: Comment out for now - database implementation needed
	// workflow, err := s.adapter.FindOne(ctx, "workflows", bson.M{"_id": req.WorkflowId}, bson.M{})
	// if err != nil {
	// 	return nil, status.Error(codes.NotFound, fmt.Sprintf("workflow not found: %v", err))
	// }

	// Generate execution ID
	executionID := uuid.New().String()

	// Convert input to map for Redis serialization
	inputMap := make(map[string]interface{})
	if req.Input != nil && req.Input.Fields != nil {
		for k, v := range req.Input.Fields {
			inputMap[k] = v.AsInterface()
		}
	}

	// Create task payload
	payload := &WorkflowTaskPayload{
		ExecutionID: executionID,
		WorkflowID:  req.WorkflowId,
		Input:       inputMap,
		Context:     req.Context,
	}

	// Create task record in database
	// TODO: Comment out for now - database implementation needed
	// taskDoc := bson.M{
	// 	"_id":         executionID,
	// 	"workflow_id": req.WorkflowId,
	// 	"input":       inputMap,
	// 	"context":     req.Context,
	// 	"status":      "QUEUED",
	// 	"created_at":  time.Now(),
	// }
	// _, err = s.adapter.InsertOne(ctx, "workflow_executions", taskDoc)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create task: %v", err))
	// }

	// Enqueue task to Redis
	_, err := s.redisClient.EnqueueWorkflow(ctx, payload)
	if err != nil {
		// Update task status to FAILED
		// TODO: Comment out for now - database implementation needed
		// s.adapter.UpdateOne(ctx, "workflow_executions",
		// 	bson.M{"_id": executionID},
		// 	bson.M{"$set": bson.M{
		// 		"status": "FAILED",
		// 		"error": err.Error(),
		// 		"updated_at": time.Now(),
		// 	}},
		// )
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to enqueue task: %v", err))
	}

	return &pb.ExecuteWorkflowResponse{
		ExecutionId: executionID,
		WorkflowId:  req.WorkflowId,
		Status:      "QUEUED",
		Message:     "Workflow execution queued successfully",
	}, nil
}

