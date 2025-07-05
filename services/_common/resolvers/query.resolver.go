package resolvers

import (
	"context"
	"fmt"
	pb "thaily/proto/common"
	"thaily/services/_common/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *CommonService) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.GenericResponse, error) {
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

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return &pb.GenericResponse{
			Success: false,
			Message: fmt.Sprintf("invalid ID format: %v", err),
		}, nil
	}
	projection := bson.M{}
	if len(req.Fields) > 0 {
		for _, field := range req.Fields {
			projection[field] = 1
		}
	}

	return s.adapter.GetById(ctx, req.EntityType, id, projection)
}

func (s *CommonService) Query(ctx context.Context, req *pb.QueryRequest) (*pb.QueryResponse, error) {
	if req.EntityType == "" {
		return &pb.QueryResponse{
			Success: false,
			Message: "entity_type is required",
		}, nil
	}

	if req.Page < 0 {
		req.Page = 1
	}

	if req.PageSize < 0 {
		req.PageSize = 10
	} else if req.PageSize > 1000 {
		req.PageSize = 1000
	}

	filter := bson.M{}

	if req.Query != "" && len(req.SearchFields) > 0 {
		orConditions := bson.A{}
		for _, field := range req.SearchFields {
			orConditions = append(orConditions, bson.M{
				field: bson.M{"$regex": req.Query, "$options": "i"},
			})
		}
		filter["$or"] = orConditions
	}

	for key, value := range req.Filters {
		filter[key] = helper.StructValueToInterface(value)
	}

	pipeline := bson.A{}
	pipelineCount := bson.A{}
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
		pipelineCount = append(pipelineCount, bson.M{"$match": filter})
	}

	for _, stage := range req.Pipeline {
		stageDoc := helper.StructToDoc(stage)
		pipeline = append(pipeline, stageDoc)
	}

	projection := bson.M{}
	if len(req.Fields) > 0 {

		for _, field := range req.Fields {
			projection[field] = 1
		}
	}

	return s.adapter.Query(ctx, req.EntityType, pipeline, projection, req.Page, req.PageSize, pipelineCount)
}
