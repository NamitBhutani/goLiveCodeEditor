package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/NamitBhutani/goLiveCodeEditor/authService/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	defer conn.Close()
	client := pb.NewAuthServiceClient(conn)

	RegisterResponse, err := client.Register(context.Background(), &pb.RegisterRequest{})
	if err != nil {
		log.Fatalf("Failed to Register %v", err)
	}
	fmt.Println(RegisterResponse.Success)

	LoginResponse, err := client.Login(context.Background(), &pb.LoginRequest{})
	if err != nil {
		log.Fatalf("Failed to login %v", err)
	}
	fmt.Println(LoginResponse.Accesstoken, LoginResponse.Refreshtoken)

	RefreshTokenResponse, err := client.RefreshToken(context.Background(), &pb.RefreshTokenRequest{})
	if err != nil {
		log.Fatalf("Failed to generate access token %v", err)
	}
	fmt.Println(RefreshTokenResponse.Accesstoken)

	LogoutResponse, err := client.Logout(context.Background(), &pb.LogoutRequest{})
	if err != nil {
		log.Fatalf("Failed to Logout %v", err)
	}
	fmt.Println(LogoutResponse.Success)
}
