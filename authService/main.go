package main

import (
	"context"
	"log"
	"net"

	pb "github.com/NamitBhutani/goLiveCodeEditor/proto"
	"google.golang.org/grpc"
)

var (
	port = ":8080"
)

type authServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

func (s *authServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{Token: "token"}, nil
}

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
