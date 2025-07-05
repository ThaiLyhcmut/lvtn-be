package common

import (
	"context"
	"fmt"
	"time"

	pb "thaily/proto/common"
	"thaily/services/_common/helper"
	"thaily/services/adapter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommonService struct {
	pb.UnimplementedCommonServiceServer
	adapter *adapter.MongoDBAdapter
}

func NewCommonService(adapter *adapter.MongoDBAdapter) *CommonService {
	return &CommonService{
		adapter: adapter,
	}
}

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
	doc["created_at"] = time.Now()
	doc["updated_at"] = time.Now()

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
		doc["created_at"] = time.Now()
		doc["updated_at"] = time.Now()
		docs[i] = doc
	}

	return s.adapter.CreateMany(ctx, req.EntityType, docs, req.Ordered)
}

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
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
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

	return s.adapter.Query(ctx, req.EntityType, pipeline, projection, req.Page, req.PageSize)
}

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
	updateDoc["updated_at"] = time.Now()

	var update bson.M
	if req.PartialUpdate {
		update = bson.M{"$set": updateDoc}
	} else {
		delete(updateDoc, "_id")
		updateDoc["updated_at"] = time.Now()
		update = bson.M{"$set": updateDoc}
	}

	return s.adapter.Update(ctx, req.EntityType, id, update)
}

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
