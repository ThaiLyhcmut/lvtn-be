package adapter

import (
	"context"
	"fmt"
	"time"

	pb "thaily/proto/common"
	"thaily/services/_common/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type MongoDBAdapter struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDBAdapter(uri, dbName string) (*MongoDBAdapter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &MongoDBAdapter{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

func (m *MongoDBAdapter) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}

func (m *MongoDBAdapter) Create(ctx context.Context, _collection string, doc bson.M) (*pb.GenericResponse, error) {
	collection := m.database.Collection(_collection)
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return &pb.GenericResponse{
			Success: false,
			Message: fmt.Sprintf("failed to create entity: %v", err),
		}, nil
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()

	doc["_id"] = result.InsertedID

	entityStruct, err := helper.DocToStruct(doc)
	if err != nil {
		return &pb.GenericResponse{
			Success: false,
			Message: fmt.Sprintf("failed to convert entity: %v", err),
		}, nil
	}

	return &pb.GenericResponse{
		Success:   true,
		Message:   "Entity created successfully",
		Id:        wrapperspb.String(id),
		Entity:    entityStruct,
		Timestamp: timestamppb.Now(),
	}, nil
}

func (m *MongoDBAdapter) CreateMany(ctx context.Context, _collection string, docs []interface{}, ordered bool) (*pb.BatchResponse, error) {
	collection := m.database.Collection(_collection)

	opts := options.InsertMany().SetOrdered(ordered)
	result, err := collection.InsertMany(ctx, docs, opts)

	if err != nil {
		if ordered {
			return &pb.BatchResponse{
				Success: false,
				Message: fmt.Sprintf("batch insert failed: %v", err),
			}, nil
		}

		bulkErr, ok := err.(mongo.BulkWriteException)
		if !ok {
			return &pb.BatchResponse{
				Success: false,
				Message: fmt.Sprintf("batch insert failed: %v", err),
			}, nil
		}

		errors := make([]*pb.ErrorDetail, len(bulkErr.WriteErrors))
		for i, we := range bulkErr.WriteErrors {
			errors[i] = &pb.ErrorDetail{
				Code:    fmt.Sprintf("%d", we.Code),
				Field:   "",
				Message: we.Message,
			}
		}

		ids := make([]string, len(result.InsertedIDs))
		for i, id := range result.InsertedIDs {
			ids[i] = id.(primitive.ObjectID).Hex()
		}

		return &pb.BatchResponse{
			Success:      false,
			Message:      "Partial batch insert completed with errors",
			Ids:          ids,
			CreatedCount: int32(len(result.InsertedIDs)),
			Errors:       errors,
		}, nil
	}

	ids := make([]string, len(result.InsertedIDs))
	entities := make([]*structpb.Struct, len(result.InsertedIDs))

	for i, id := range result.InsertedIDs {
		ids[i] = id.(primitive.ObjectID).Hex()

		var entity bson.M
		err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
		if err == nil {
			if entityStruct, err := helper.DocToStruct(entity); err == nil {
				entities[i] = entityStruct
			}
		}
	}

	return &pb.BatchResponse{
		Success:      true,
		Message:      "Batch insert completed successfully",
		Ids:          ids,
		CreatedCount: int32(len(result.InsertedIDs)),
		Entities:     entities,
	}, nil
}

func (m *MongoDBAdapter) GetById(ctx context.Context, _collection string, _id primitive.ObjectID, projection bson.M) (*pb.GenericResponse, error) {
	collection := m.database.Collection(_collection)

	opts := options.FindOne()
	if len(projection) > 0 {
		opts.SetProjection(projection)
	}
	var entity bson.M
	err := collection.FindOne(ctx, bson.M{"_id": _id}, opts).Decode(&entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.GenericResponse{
				Success: false,
				Message: "Entity not found",
			}, nil
		}
		return &pb.GenericResponse{
			Success: false,
			Message: fmt.Sprintf("failed to get entity: %v", err),
		}, nil
	}

	entityStruct, err := helper.DocToStruct(entity)
	if err != nil {
		return &pb.GenericResponse{
			Success: false,
			Message: fmt.Sprintf("failed to convert entity: %v", err),
		}, nil
	}

	return &pb.GenericResponse{
		Success:   true,
		Message:   "Entity retrieved successfully",
		Id:        wrapperspb.String(_id.Hex()),
		Entity:    entityStruct,
		Timestamp: timestamppb.Now(),
	}, nil
}

func (m *MongoDBAdapter) Query(ctx context.Context, _collection string, pipeline bson.A, projection bson.M, page int32, pagesize int32, pipelineCount bson.A) (*pb.QueryResponse, error) {
	collection := m.database.Collection(_collection)
	type countResult struct {
		total int64
		err   error
	}
	type dataResult struct {
		results []bson.M
		err     error
	}

	countChan := make(chan countResult, 1)
	dataChan := make(chan dataResult, 1)

	// Goroutine cho count query
	go func() {
		countPipeline := append(pipelineCount, bson.M{"$count": "total"})
		countCursor, err := collection.Aggregate(ctx, countPipeline)
		if err != nil {
			countChan <- countResult{err: err}
			return
		}
		defer countCursor.Close(ctx)

		var countResults []bson.M
		if err := countCursor.All(ctx, &countResults); err != nil {
			countChan <- countResult{err: err}
			return
		}

		totalItems := int64(0)
		if len(countResults) > 0 {
			if total, ok := countResults[0]["total"].(int32); ok {
				totalItems = int64(total)
			} else if total, ok := countResults[0]["total"].(int64); ok {
				totalItems = total
			}
		}

		countChan <- countResult{total: totalItems}
	}()

	// Goroutine cho data query
	go func() {
		// Add pagination
		if page > 0 && pagesize > 0 {
			skip := (page - 1) * pagesize
			pipeline = append(pipeline, bson.M{"$skip": skip})
			pipeline = append(pipeline, bson.M{"$limit": pagesize})
		}

		// Add projection
		if len(projection) > 0 {
			pipeline = append(pipeline, bson.M{"$project": projection})
		}

		cursor, err := collection.Aggregate(ctx, pipeline)
		if err != nil {
			dataChan <- dataResult{err: err}
			return
		}
		defer cursor.Close(ctx)

		var results []bson.M
		if err := cursor.All(ctx, &results); err != nil {
			dataChan <- dataResult{err: err}
			return
		}

		dataChan <- dataResult{results: results}
	}()

	// Đợi cả 2 kết quả
	countRes := <-countChan
	dataRes := <-dataChan

	// Xử lý lỗi
	if countRes.err != nil {
		return &pb.QueryResponse{
			Success: false,
			Message: fmt.Sprintf("failed to count entities: %v", countRes.err),
		}, nil
	}

	if dataRes.err != nil {
		return &pb.QueryResponse{
			Success: false,
			Message: fmt.Sprintf("failed to query entities: %v", dataRes.err),
		}, nil
	}

	// Convert results
	entities := make([]*structpb.Struct, len(dataRes.results))
	for i, result := range dataRes.results {
		if entityStruct, err := helper.DocToStruct(result); err == nil {
			entities[i] = entityStruct
		}
	}

	// Calculate pagination
	totalPages := int32(0)
	if pagesize > 0 {
		totalPages = int32((countRes.total + int64(pagesize) - 1) / int64(pagesize))
	}

	return &pb.QueryResponse{
		Success:  true,
		Message:  "Query executed successfully",
		Entities: entities,
		Pagination: &pb.Pagination{
			CurrentPage: page,
			PageSize:    pagesize,
			TotalPages:  totalPages,
			TotalItems:  countRes.total,
		},
	}, nil
}

func (m *MongoDBAdapter) Update(ctx context.Context, _collection string, _id primitive.ObjectID, update bson.M) (*pb.GenericResponse, error) {
	collection := m.database.Collection(_collection)
	var updatedEntity bson.M
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": _id}, update, opts).Decode(&updatedEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.GenericResponse{
				Success: false,
				Message: "Entity not found",
			}, nil
		}
		return &pb.GenericResponse{
			Success: false,
			Message: fmt.Sprintf("failed to update entity: %v", err),
		}, nil
	}

	entityStruct, err := helper.DocToStruct(updatedEntity)
	if err != nil {
		return &pb.GenericResponse{
			Success: false,
			Message: fmt.Sprintf("failed to convert entity: %v", err),
		}, nil
	}

	return &pb.GenericResponse{
		Success:   true,
		Message:   "Entity updated successfully",
		Id:        wrapperspb.String(_id.Hex()),
		Entity:    entityStruct,
		Timestamp: timestamppb.Now(),
	}, nil
}

func (m *MongoDBAdapter) Delete(ctx context.Context, _collection string, _id primitive.ObjectID) (*pb.DeleteResponse, error) {
	collection := m.database.Collection(_collection)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": _id})
	if err != nil {
		return &pb.DeleteResponse{
			Success: false,
			Message: fmt.Sprintf("failed to delete entity: %v", err),
		}, nil
	}

	if result.DeletedCount == 0 {
		return &pb.DeleteResponse{
			Success: false,
			Message: "Entity not found",
		}, nil
	}

	return &pb.DeleteResponse{
		Success:      true,
		Message:      "Entity deleted successfully",
		DeletedCount: int32(result.DeletedCount),
	}, nil
}

func (m *MongoDBAdapter) DeleteMany(ctx context.Context, _collection string, ids []primitive.ObjectID, pipeline bson.A) (*pb.DeleteManyResponse, error) {
	collection := m.database.Collection(_collection)
	filter := bson.M{}
	filter["_id"] = bson.M{"$in": ids}

	if len(pipeline) > 0 {

		pipeline = append(pipeline, bson.M{"$match": filter})

		cursor, err := collection.Aggregate(ctx, pipeline)
		if err != nil {
			return &pb.DeleteManyResponse{
				Success: false,
				Message: fmt.Sprintf("failed to execute pipeline: %v", err),
			}, nil
		}
		defer cursor.Close(ctx)

		var docsToDelete []bson.M
		if err := cursor.All(ctx, &docsToDelete); err != nil {
			return &pb.DeleteManyResponse{
				Success: false,
				Message: fmt.Sprintf("failed to get documents: %v", err),
			}, nil
		}

		ids := make([]primitive.ObjectID, len(docsToDelete))
		for i, doc := range docsToDelete {
			ids[i] = doc["_id"].(primitive.ObjectID)
		}

		filter = bson.M{"_id": bson.M{"$in": ids}}
	}

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return &pb.DeleteManyResponse{
			Success: false,
			Message: fmt.Sprintf("failed to delete entities: %v", err),
		}, nil
	}

	return &pb.DeleteManyResponse{
		Success:      true,
		Message:      "Entities deleted successfully",
		DeletedCount: int32(result.DeletedCount),
		FailedIds:    []string{},
	}, nil
}

func (m *MongoDBAdapter) Aggregate(ctx context.Context, _collection string, allowDiskUse bool, maxTimeMs int32, pipeline bson.A) (*pb.AggregateResponse, error) {
	collection := m.database.Collection(_collection)
	opts := options.Aggregate()
	if allowDiskUse {
		opts.SetAllowDiskUse(true)
	}
	if maxTimeMs > 0 {
		opts.SetMaxTime(time.Duration(maxTimeMs) * time.Millisecond)
	}

	startTime := time.Now()
	cursor, err := collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return &pb.AggregateResponse{
			Success: false,
			Message: fmt.Sprintf("aggregation failed: %v", err),
		}, nil
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return &pb.AggregateResponse{
			Success: false,
			Message: fmt.Sprintf("failed to get results: %v", err),
		}, nil
	}

	resultsStruct := make([]*structpb.Struct, len(results))
	for i, result := range results {
		if resultStruct, err := helper.DocToStruct(result); err == nil {
			resultsStruct[i] = resultStruct
		}
	}

	executionTime := time.Since(startTime).Milliseconds()

	return &pb.AggregateResponse{
		Success:         true,
		Message:         "Aggregation completed successfully",
		Results:         resultsStruct,
		ExecutionTimeMs: executionTime,
	}, nil
}
