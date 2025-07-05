package resolvers

import (
	"context"
	pb "thaily/proto/common"
	"thaily/services/_common/helper"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *CommonService) Aggregate(ctx context.Context, req *pb.AggregateRequest) (*pb.AggregateResponse, error) {
	if req.EntityType == "" {
		return &pb.AggregateResponse{
			Success: false,
			Message: "entity_type is required",
		}, nil
	}

	if len(req.Pipeline) == 0 {
		return &pb.AggregateResponse{
			Success: false,
			Message: "pipeline cannot be empty",
		}, nil
	}

	if req.MaxTimeMs > 300000 {
		req.MaxTimeMs = 300000
	}

	pipeline := bson.A{}
	for _, stage := range req.Pipeline {
		stageDoc := helper.StructToDoc(stage)
		pipeline = append(pipeline, stageDoc)
	}

	return s.adapter.Aggregate(ctx, req.EntityType, req.AllowDiskUse, req.MaxTimeMs, pipeline)
}
