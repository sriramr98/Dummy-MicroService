package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sriramr98/todo_auth_service/services"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
)

type AuthGrpcServer struct {
	authService services.AuthService
}

func GetGRPCServer(authService services.AuthService) (*grpc.Server, net.Listener) {
	grpcAddr := fmt.Sprintf("localhost:%s", os.Getenv("GRPC_PORT"))
	log.Printf("Starting GRPC Server at %s", grpcAddr)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	services.RegisterAuthServiceServer(grpcServer, AuthGrpcServer{authService: authService})

	return grpcServer, lis
}

func (asg AuthGrpcServer) GetUserInfo(ctx context.Context, ur *services.UserRequest) (*services.UserInfo, error) {
	log.Println("GetUserInfo hit..")
	token := ur.GetToken()

	tokenRsp, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Println(err)
		// handle err
		return nil, err
	}

	// validate the essential claims
	if !tokenRsp.Valid {
		// handle invalid tokebn
		return nil, errors.New("invalid token")
	}

	iss, err := tokenRsp.Claims.GetSubject()
	if err != nil {
		log.Println(err)
		// handle error
		return nil, err
	}
	userId, err := strconv.Atoi(iss)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	user, err := asg.authService.GetUserInfo(int64(userId))
	if err != nil {
		log.Println(err)
		//TODO: Handle error
		return &services.UserInfo{}, err
	}

	return &services.UserInfo{
		Email:   user.Email,
		Blocked: user.Blocked,
		UserId:  int64(userId),
	}, nil
}
