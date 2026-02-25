package main

import (
	"log"
	"net"

	"github.com/mohammad-kh1/distributed-gateway/api/proto"
	"github.com/mohammad-kh1/distributed-gateway/internal/auth"
	dbAuth "github.com/mohammad-kh1/distributed-gateway/internal/db"

	"google.golang.org/grpc"
)

func main() {
	// connect to postgres database
	dsn := "postgres://postgres:password@localhost:5432/auth_db"
	repo := dbAuth.NewRepository(dsn)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v \n", err)
	}

	s := grpc.NewServer()

	proto.RegisterAuthServiceServer(s, auth.NewAuthServer(repo))

	log.Println("Auth Service os runnin on port 50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
