package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/NamitBhutani/goLiveCodeEditor/proto"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
)

var (
	port          = ":8080"
	secret        = []byte("secret")
	accessExpire  = time.Hour      // Access token expiration time
	refreshExpire = time.Hour * 24 // Refresh token expiration time
)

type authServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

func generateToken(username string, expireTime time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(expireTime).Unix()
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateToken(tokenString string, expectedExpire time.Duration) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Check token expiration
	expiration := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expiration) {
		return nil, fmt.Errorf("token has expired")
	}

	// Check token expiration time against the expected expire time
	if expectedExpire > 0 && time.Until(expiration) > expectedExpire {
		return nil, fmt.Errorf("token expiration time too long")
	}

	return token, nil
}

func (s *authServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Generate access token
	accessToken, err := generateToken(req.Username, accessExpire)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := generateToken(req.Username, refreshExpire)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Accesstoken:  accessToken,
		Refreshtoken: refreshToken,
	}, nil
}

func (s *authServiceServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// Validate the refresh token
	_, err := validateToken(req.Refreshtoken, refreshExpire)
	if err != nil {
		return nil, err
	}

	// Generate a new access token
	newAccessToken, err := generateToken(req.Username, accessExpire)
	if err != nil {
		return nil, err
	}

	return &pb.RefreshTokenResponse{
		Accesstoken: newAccessToken,
	}, nil
}

func (s *authServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Perform registration logic here
	return &pb.RegisterResponse{Success: true}, nil
}

func (s *authServiceServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	token, err := validateToken(req.Token, accessExpire)
	if err != nil {
		return nil, err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return &pb.LogoutResponse{Success: true}, nil
}

// server setup
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &authServiceServer{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
