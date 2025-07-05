package resolvers

import (
	"context"
	pb "thaily/proto/common"
	"thaily/services/_common/helper"
	"time"
)

func (s *CommonService) Create(ctx context.Context, req *pb.GenericRequest) (*pb.GenericResponse, error) {
	if req.EntityType == "" {
		return &pb.GenericResponse{
			Success: false,
			Message: "entity_type is required",
		}, nil
	}

	if req.Data == nil {
		return &pb.GenericResponse{
			Success: false,
			Message: "data is required",
		}, nil
	}

	doc := helper.StructToDoc(req.Data)
	now := time.Now()
	doc["createdAt"] = now
	doc["updatedAt"] = now

	return s.adapter.Create(ctx, req.EntityType, doc)
}

func (s *CommonService) CreateMany(ctx context.Context, req *pb.BatchRequest) (*pb.BatchResponse, error) {
	if req.EntityType == "" {
		return &pb.BatchResponse{
			Success: false,
			Message: "entity_type is required",
		}, nil
	}

	if len(req.Entities) == 0 {
		return &pb.BatchResponse{
			Success: false,
			Message: "entities array cannot be empty",
		}, nil
	}

	docs := make([]interface{}, len(req.Entities))
	for i, entity := range req.Entities {
		doc := helper.StructToDoc(entity)
		now := time.Now()
		doc["createdAt"] = now
		doc["updatedAt"] = now
		docs[i] = doc
	}

	return s.adapter.CreateMany(ctx, req.EntityType, docs, req.Ordered)
}
