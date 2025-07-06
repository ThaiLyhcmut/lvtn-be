package adapter

import (
	"context"

	pb "thaily/proto/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DatabaseAdapter interface {
	Create(ctx context.Context, _collection string, doc bson.M) (*pb.GenericResponse, error)
	CreateMany(ctx context.Context, _collection string, docs []interface{}, ordered bool) (*pb.BatchResponse, error)
	GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.GenericResponse, error)
	Query(ctx context.Context, _collection string, pipeline bson.A, projection bson.M, page int32, pagesize int32, pipelineCount bson.A)
	FindOne(ctx context.Context, _collection string, conditions bson.M, projection bson.M) (*pb.GenericResponse, error)
	Update(ctx context.Context, _collection string, _id primitive.ObjectID, update bson.M) (*pb.GenericResponse, error)
	Delete(ctx context.Context, _collection string, _id primitive.ObjectID) (*pb.DeleteResponse, error)
	DeleteMany(ctx context.Context, _collection string, ids []primitive.ObjectID, pipeline bson.A) (*pb.DeleteManyResponse, error)
	Aggregate(ctx context.Context, _collection string, allowDiskUse bool, maxTimeMs int32, pipeline bson.A) (*pb.AggregateResponse, error)
	Close() error
}
