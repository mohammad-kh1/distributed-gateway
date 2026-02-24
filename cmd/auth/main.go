package main

import (
	"log"
	"net"

	"github.com/mohammad-kh1/distributed-gateway/api/proto"
	"github.com/mohammad-kh1/distributed-gateway/internal/auth"
	"google.golang.org/grpc"
)

func main(){
	lis , err := net.Listen("tcp" , ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v \n",err)
	}

	s := grpc.NewServer()

	proto.RegisterAuthServiceServer(s , &auth.Server{})

	log.Println("Auth Service os runnin on port 50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n",err)
	}
}
