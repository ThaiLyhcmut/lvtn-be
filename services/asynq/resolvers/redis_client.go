package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

const (
	// Task types
	TypeWorkflowExecution = "workflow:execute"
	
	// Queue names
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

// WorkflowTaskPayload represents the payload for workflow execution tasks
type WorkflowTaskPayload struct {
	ExecutionID string            `json:"execution_id"`
	WorkflowID  string            `json:"workflow_id"`
	Input       map[string]interface{} `json:"input"`
	Context     map[string]string `json:"context"`
}

// RedisTaskClient wraps asynq client for task management
type RedisTaskClient struct {
	client    *asynq.Client
	inspector *asynq.Inspector
}

// NewRedisTaskClient creates a new Redis task client
func NewRedisTaskClient(redisAddr string, redisDB int) *RedisTaskClient {
	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
		DB:   redisDB,
	}
	
	return &RedisTaskClient{
		client:    asynq.NewClient(redisOpt),
		inspector: asynq.NewInspector(redisOpt),
	}
}

// EnqueueWorkflow enqueues a workflow execution task
func (r *RedisTaskClient) EnqueueWorkflow(ctx context.Context, payload *WorkflowTaskPayload, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	task := asynq.NewTask(TypeWorkflowExecution, data)
	
	// Default options
	defaultOpts := []asynq.Option{
		asynq.Queue(QueueDefault),
		asynq.MaxRetry(3),
		asynq.Timeout(30 * time.Minute),
		asynq.Retention(24 * time.Hour),
	}
	
	// Merge with provided options
	opts = append(defaultOpts, opts...)
	
	info, err := r.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue task: %w", err)
	}
	
	return info, nil
}

// GetTaskInfo retrieves task information by ID
func (r *RedisTaskClient) GetTaskInfo(ctx context.Context, taskID string) (*asynq.TaskInfo, error) {
	// Search in all queues
	queues := []string{QueueCritical, QueueDefault, QueueLow}
	
	for _, queue := range queues {
		info, err := r.inspector.GetTaskInfo(queue, taskID)
		if err == nil {
			return info, nil
		}
	}
	
	return nil, fmt.Errorf("task not found: %s", taskID)
}

// Close closes the Redis client connection
func (r *RedisTaskClient) Close() error {
	return r.client.Close()
}