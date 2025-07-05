package resolvers

import (
	"context"
	"fmt"
	pb "thaily/proto/common"
	"thaily/services/_common/helper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *CommonService) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.GenericResponse, error) {
	if req.EntityType == "" {
		return &pb.GenericResponse{
			Success: false,
			Message: "entity_type is required",
		}, nil
	}

	if req.Id == "" {
		return &pb.GenericResponse{
			Success: false,
			Message: "id is required",
		}, nil
	}

	if req.Data == nil {
		return &pb.GenericResponse{
			Success: false,
			Message: "data is required",
		}, nil
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return &pb.GenericResponse{
			Success: false,
			Message: fmt.Sprintf("invalid ID format: %v", err),
		}, nil
	}

	updateDoc := helper.StructToDoc(req.Data)
	updateDoc["updatedAt"] = time.Now()

	var update bson.M
	if req.PartialUpdate {
		update = bson.M{"$set": updateDoc}
	} else {
		delete(updateDoc, "_id")
		update = bson.M{"$set": updateDoc}
	}

	return s.adapter.Update(ctx, req.EntityType, id, update)
}
