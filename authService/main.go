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
	port       = ":8080"
	secret     = []byte("secret")
	expireTime = time.Hour * 24
)

type authServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

func (s *authServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	//token generation
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = req.Username
	claims["exp"] = time.Now().Add(expireTime).Unix()
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{Token: tokenString}, nil
}

func validateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func (s *authServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	//db calls
	return &pb.RegisterResponse{Success: true}, nil
}

func (s *authServiceServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	token, err := validateToken(req.Token)
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

	grpcServer := grpc.NewServer()                                 //init new server
	pb.RegisterAuthServiceServer(grpcServer, &authServiceServer{}) //register the service
	if err := grpcServer.Serve(lis); err != nil {                  //serve the service
		log.Fatalf("failed to serve: %v", err)
	}
}
