package utils

import (
	"context"
	"fmt"

	pb "thaily/proto/asynq"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	ID           string                 `bson:"id"`
	Name         string                 `bson:"name"`
	ClientName   string                 `bson:"client_name"`
	Method       string                 `bson:"method"`
	InputMapping map[string]interface{} `bson:"input_mapping"`
	DependsOn    []string               `bson:"depends_on"`
	Retry        int                    `bson:"retry"`
	Timeout      int                    `bson:"timeout"`
}

// Workflow represents a workflow definition
type Workflow struct {
	ID          string         `bson:"_id"`
	Name        string         `bson:"name"`
	Description string         `bson:"description"`
	Steps       []WorkflowStep `bson:"steps"`
	CreatedAt   int64          `bson:"created_at"`
	UpdatedAt   int64          `bson:"updated_at"`
}

// ExecuteWorkflowSteps executes workflow steps with the given input
func ExecuteWorkflowSteps(ctx context.Context, workflow *Workflow, input *pb.DynamicInput, clients map[string]*grpc.ClientConn) (*structpb.Struct, error) {
	// Store results from each step
	stepResults := make(map[string]*structpb.Value)

	// Add initial input to results
	if input != nil && input.Fields != nil {
		for k, v := range input.Fields {
			stepResults[fmt.Sprintf("input.%s", k)] = v
		}
	}

	// Execute each step in order
	for _, step := range workflow.Steps {
		// Check dependencies
		for _, dep := range step.DependsOn {
			if _, ok := stepResults[dep]; !ok {
				return nil, fmt.Errorf("dependency %s not satisfied for step %s", dep, step.ID)
			}
		}

		// Get client for this step
		client, ok := clients[step.ClientName]
		if !ok {
			return nil, fmt.Errorf("client %s not found for step %s", step.ClientName, step.ID)
		}

		// Prepare input for this step
		stepInput := prepareStepInput(step.InputMapping, stepResults)

		// Execute step (placeholder for actual gRPC call)
		// TODO: Implement actual gRPC method invocation
		result := &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: fmt.Sprintf("Result from step %s", step.ID),
			},
		}

		// Store result
		stepResults[step.ID] = result
	}

	// Prepare final output
	output := &structpb.Struct{
		Fields: make(map[string]*structpb.Value),
	}
	for k, v := range stepResults {
		output.Fields[k] = v
	}

	return output, nil
}

// prepareStepInput maps input data based on the step's input mapping configuration
func prepareStepInput(inputMapping map[string]interface{}, stepResults map[string]*structpb.Value) *structpb.Struct {
	input := &structpb.Struct{
		Fields: make(map[string]*structpb.Value),
	}

	for targetKey, sourceKey := range inputMapping {
		if sourceKeyStr, ok := sourceKey.(string); ok {
			if value, exists := stepResults[sourceKeyStr]; exists {
				input.Fields[targetKey] = value
			}
		}
	}

	return input
}
