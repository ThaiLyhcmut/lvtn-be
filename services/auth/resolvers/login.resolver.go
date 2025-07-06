package resolver

import (
	"context"
	"fmt"
	pb "thaily/proto/auth"
	"thaily/services/_common/helper"
	"thaily/services/auth/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return &pb.LoginResponse{
			Success: false,
			Message: "Email and password are required",
		}, nil
	}

	passwordHash, ok := utils.HashPassword(req.Password)
	if ok != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: "Invalid password",
		}, nil
	}

	conditions := bson.M{
		"email":    req.Email,
		"password": passwordHash,
	}

	// push limit

	data, err := s.adapter.FindOne(ctx, "users", conditions, bson.M{})

	user := helper.StructToDoc(data.Entity)

	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: "Error Login with email password",
		}, nil
	}

	// Check if user is active
	if status, ok := user["status"].(string); ok && status != "active" {
		return &pb.LoginResponse{
			Success: false,
			Message: "Your account is not active",
		}, nil
	}

	// Extract user info
	userID := data.Id.GetValue()
	email := user["email"].(string)
	fullName, _ := user["fullName"].(string)
	if fullName == "" {
		fullName, _ = user["name"].(string)
	}
	roles, _ := user["roles"].(string)
	if roles == "" {
		roles = "user"
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(userID, email, fullName, roles)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate refresh token")
	}

	// // Store refresh token in database (optional)
	// refreshCollection := s.adapter.GetDatabase().Collection("refresh_tokens")
	// refreshCollection.InsertOne(ctx, bson.M{
	// 	"userId":    user["_id"],
	// 	"token":     refreshToken,
	// 	"createdAt": time.Now(),
	// 	"expiresAt": time.Now().Add(7 * 24 * time.Hour),
	// })
	fmt.Println(user["createdAt"])
	createdAt, _ := user["createdAt"].(time.Time)
	updatedAt, _ := user["updatedAt"].(string)
	fmt.Print(createdAt, updatedAt)
	return &pb.LoginResponse{
		Success:      true,
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &pb.User{
			Id:        userID,
			Code:      fmt.Sprintf("USR%s", userID[:6]),
			Email:     email,
			FullName:  fullName,
			Roles:     roles,
			AvatarUrl: "",
			CreatedAt: 0,
			UpdatedAt: 0,
		},
	}, nil
}

func (s *AuthService) GoogleLogin(ctx context.Context, req *pb.GoogleLoginRequest) (*pb.LoginResponse, error) {
	// TODO: Implement Google OAuth verification
	// 1. Verify Google ID token
	// 2. Extract user info from token
	// 3. Create or update user in database
	// 4. Generate JWT tokens

	return &pb.LoginResponse{
		Success: false,
		Message: "Google login not implemented yet",
	}, nil
}
