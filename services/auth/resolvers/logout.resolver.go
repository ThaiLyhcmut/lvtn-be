package resolver

import (
	"context"
	pb "thaily/proto/auth"
)

func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if req.AccessToken == "" {
		return &pb.LogoutResponse{
			Success: false,
			Message: "Access token is required",
		}, nil
	}

	// Verify token
	_, err := s.jwtManager.VerifyToken(req.AccessToken)
	if err != nil {
		return &pb.LogoutResponse{
			Success: false,
			Message: "Invalid access token",
		}, nil
	}

	// Remove refresh tokens for this user
	// refreshCollection := s.adapter.GetDatabase().Collection("refresh_tokens")
	// userObjID, _ := primitive.ObjectIDFromHex(claims.UserID)
	// _, err = refreshCollection.DeleteMany(ctx, bson.M{"userId": userObjID})
	// if err != nil {
	// 	return &pb.LogoutResponse{
	// 		Success: false,
	// 		Message: "Failed to logout",
	// 	}, nil
	// }

	return &pb.LogoutResponse{
		Success: true,
		Message: "Logout successful",
	}, nil
}
