package resolver

import (
	pb "thaily/proto/auth"
	"thaily/services/adapter"
	"thaily/services/auth/utils"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	adapter    *adapter.MongoDBAdapter
	jwtManager *utils.JWTManager
}

func NewAuthService(adapter *adapter.MongoDBAdapter, jwtSecret string) *AuthService {
	return &AuthService{
		adapter:    adapter,
		jwtManager: utils.NewJWTManager(jwtSecret),
	}
}
