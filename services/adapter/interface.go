package adapter

import (
	"context"

	pb "thaily/proto/common"
)

type DatabaseAdapter interface {
	Create(ctx context.Context, req *pb.GenericRequest) (*pb.GenericResponse, error)
	CreateMany(ctx context.Context, req *pb.BatchRequest) (*pb.BatchResponse, error)
	GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.GenericResponse, error)
	Query(ctx context.Context, req *pb.QueryRequest) (*pb.QueryResponse, error)
	Update(ctx context.Context, req *pb.UpdateRequest) (*pb.GenericResponse, error)
	Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error)
	DeleteMany(ctx context.Context, req *pb.DeleteManyRequest) (*pb.DeleteManyResponse, error)
	Aggregate(ctx context.Context, req *pb.AggregateRequest) (*pb.AggregateResponse, error)
	Close() error
}
