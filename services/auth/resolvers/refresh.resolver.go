package resolver

import (
	"context"
	pb "thaily/proto/auth"
)

func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	return nil, nil
	// if req.RefreshToken == "" {
	// 	return &pb.RefreshTokenResponse{
	// 		Success: false,
	// 		Message: "Refresh token is required",
	// 	}, nil
	// }

	// // Verify refresh token
	// claims, err := s.jwtManager.VerifyToken(req.RefreshToken)
	// if err != nil {
	// 	return &pb.RefreshTokenResponse{
	// 		Success: false,
	// 		Message: "Invalid refresh token",
	// 	}, nil
	// }

	// // Check if refresh token exists in database
	// refreshCollection := s.adapter.GetDatabase().Collection("refresh_tokens")
	// userObjID, _ := primitive.ObjectIDFromHex(claims.UserID)
	// var tokenDoc bson.M
	// err = refreshCollection.FindOne(ctx, bson.M{
	// 	"userId": userObjID,
	// 	"token":  req.RefreshToken,
	// }).Decode(&tokenDoc)
	// if err != nil {
	// 	return &pb.RefreshTokenResponse{
	// 		Success: false,
	// 		Message: "Invalid refresh token",
	// 	}, nil
	// }

	// // Get user info
	// collection := s.adapter.GetDatabase().Collection("users")
	// var user bson.M
	// err = collection.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&user)
	// if err != nil {
	// 	return &pb.RefreshTokenResponse{
	// 		Success: false,
	// 		Message: "User not found",
	// 	}, nil
	// }

	// // Extract user info
	// email := user["email"].(string)
	// fullName, _ := user["fullName"].(string)
	// if fullName == "" {
	// 	fullName, _ = user["name"].(string)
	// }
	// roles, _ := user["roles"].(string)
	// if roles == "" {
	// 	roles = "user"
	// }

	// // Generate new tokens
	// newAccessToken, err := s.jwtManager.GenerateAccessToken(claims.UserID, email, fullName, roles)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, "Failed to generate access token")
	// }

	// newRefreshToken, err := s.jwtManager.GenerateRefreshToken(claims.UserID)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, "Failed to generate refresh token")
	// }

	// // Delete old refresh token and insert new one
	// refreshCollection.DeleteOne(ctx, bson.M{"token": req.RefreshToken})
	// refreshCollection.InsertOne(ctx, bson.M{
	// 	"userId":    userObjID,
	// 	"token":     newRefreshToken,
	// 	"createdAt": time.Now(),
	// 	"expiresAt": time.Now().Add(7 * 24 * time.Hour),
	// })

	// return &pb.RefreshTokenResponse{
	// 	Success:      true,
	// 	Message:      "Token refreshed successfully",
	// 	AccessToken:  newAccessToken,
	// 	RefreshToken: newRefreshToken,
	// }, nil
}
