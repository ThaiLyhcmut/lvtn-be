package resolvers

import (
	"context"
	"fmt"
	pb "thaily/proto/common"
	"thaily/services/_common/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *CommonService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if req.EntityType == "" {
		return &pb.DeleteResponse{
			Success: false,
			Message: "entity_type is required",
		}, nil
	}

	if req.Id == "" {
		return &pb.DeleteResponse{
			Success: false,
			Message: "id is required",
		}, nil
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return &pb.DeleteResponse{
			Success: false,
			Message: fmt.Sprintf("invalid ID format: %v", err),
		}, nil
	}

	return s.adapter.Delete(ctx, req.EntityType, id)
}

func (s *CommonService) DeleteMany(ctx context.Context, req *pb.DeleteManyRequest) (*pb.DeleteManyResponse, error) {
	if req.EntityType == "" {
		return &pb.DeleteManyResponse{
			Success: false,
			Message: "entity_type is required",
		}, nil
	}

	if len(req.Ids) == 0 && len(req.Pipeline) == 0 {
		return &pb.DeleteManyResponse{
			Success: false,
			Message: "either ids or pipeline must be provided",
		}, nil
	}

	objectIDs := make([]primitive.ObjectID, 0, len(req.Ids))
	if len(req.Ids) > 0 {
		objectIDs := make([]primitive.ObjectID, 0, len(req.Ids))
		failedIDs := make([]string, 0)

		for _, idStr := range req.Ids {
			id, err := primitive.ObjectIDFromHex(idStr)
			if err != nil {
				failedIDs = append(failedIDs, idStr)
				continue
			}
			objectIDs = append(objectIDs, id)
		}

		if len(failedIDs) > 0 && len(objectIDs) == 0 {
			return &pb.DeleteManyResponse{
				Success:   false,
				Message:   "All provided IDs are invalid",
				FailedIds: failedIDs,
			}, nil
		}
	}
	pipeline := bson.A{}
	for _, stage := range req.Pipeline {
		stageDoc := helper.StructToDoc(stage)
		pipeline = append(pipeline, stageDoc)
	}

	return s.adapter.DeleteMany(ctx, req.EntityType, objectIDs, pipeline)
}
