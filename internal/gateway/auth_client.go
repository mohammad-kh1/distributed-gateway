package gateway

import (
	"context"
	"log"

	"github.com/mohammad-kh1/distributed-gateway/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client proto.AuthServiceClient
}

func NewAuthClient(addr string) *AuthClient {
	//connec to grpc server

	conn , err := grpc.Dial(addr , grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to auth service : %v",err)
	}

	return &AuthClient{
		client:proto.NewAuthServiceClient(conn),
	}

}

func (ac *AuthClient) Authenticate(token string) (*proto.VerifyResponse , error){
	return ac.client.VerifyToken(context.Background() , &proto.VerifyRequest{
		Token:token,
	})
}
