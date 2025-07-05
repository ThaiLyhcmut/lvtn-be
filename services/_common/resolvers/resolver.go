package resolvers

import (
	pb "thaily/proto/common"
	"thaily/services/adapter"
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
