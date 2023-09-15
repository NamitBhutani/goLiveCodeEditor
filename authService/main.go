package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/NamitBhutani/goLiveCodeEditor/database"
	pb "github.com/NamitBhutani/goLiveCodeEditor/proto"
	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	port          = ":8080"
	secret        = []byte(os.Getenv("SECRET_KEY"))
	accessExpire  = time.Hour      // Access token expiration time
	refreshExpire = time.Hour * 24 // Refresh token expiration time
	db            *gorm.DB
)

type authServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

// db conn
func init() {
	var err error
	db, err = database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(&database.User{}, &database.TokenBlacklist{})
	if err != nil {
		log.Fatalf("failed to automigrate tables: %v", err)
	}
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
	var count int64
	err := db.Model(&database.TokenBlacklist{}).Where("token = ?", tokenString).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("token is invalid")
	}
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
	var user database.User
	result := db.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, result.Error
	}

	// Compare the hashed password from the database with the provided password.
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}
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
	user.RefreshToken = refreshToken
	db.Save(&user)

	return &pb.LoginResponse{
		Accesstoken:  accessToken,
		Refreshtoken: refreshToken,
	}, nil
}

func (s *authServiceServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// Validate the refresh token
	token, err := validateToken(req.Refreshtoken, refreshExpire)
	if err != nil {
		return nil, err
	}

	// Extract the username from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username in token claims")
	}

	// Generate a new access token
	accessToken, err := generateToken(username, accessExpire)
	if err != nil {
		return nil, err
	}

	return &pb.RefreshTokenResponse{
		Accesstoken: accessToken,
	}, nil
}

func (s *authServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	username := req.Username
	password := req.Password
	if username == "" {
		return nil, fmt.Errorf("username is empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password is empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	var existingUser database.User
	result := db.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		return nil, fmt.Errorf("username is already taken")
	}
	newUser := database.User{
		Username: username,
		Password: string(hashedPassword),
	}

	err = db.Create(&newUser).Error
	if err != nil {
		return &pb.RegisterResponse{Success: false}, err
	}
	return &pb.RegisterResponse{Success: true}, nil
}

func (s *authServiceServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	token, err := validateToken(req.Refreshtoken, refreshExpire)
	if err != nil {
		return nil, err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	tokenString := token.Raw
	err = db.Create(&database.TokenBlacklist{Token: tokenString}).Error
	if err != nil {
		return nil, err
	}
	return &pb.LogoutResponse{Success: true}, nil
}

// server setup
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go runProxyServer()
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &authServiceServer{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
func runProxyServer() {
	grpcMuxHandler := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := pb.RegisterAuthServiceHandlerServer(ctx, grpcMuxHandler, &authServiceServer{})
	if err != nil {
		log.Fatalf("failed to register REST gateway: %v", err)
	}
	muxServe := http.NewServeMux()
	muxServe.Handle("/", grpcMuxHandler)
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = http.Serve(lis, muxServe)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	//for future ref: grpcMuxHandler is the handler made by grpc-gateway,
	//which is then passed to http.Serve to handle and serve it as http request
}
