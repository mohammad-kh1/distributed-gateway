package auth

import (
	"context"
	"log"

	"github.com/mohammad-kh1/distributed-gateway/api/proto"
	auth "github.com/mohammad-kh1/distributed-gateway/internal/db"
)

type Server struct {
	proto.UnimplementedAuthServiceServer
	repo *auth.Repository
}

func NewAuthServer(repo *auth.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) VerifyToken(ctx context.Context, req *proto.VerifyRequest) (*proto.VerifyResponse, error) {
	user, err := s.repo.GetUserByToken(ctx, req.Token)
	if err != nil {
		log.Println(err)
		return &proto.VerifyResponse{Authorized: false}, nil
	}
	log.Printf(" User found: %s with limit %d", user.UserID, user.RateLimit)
	return &proto.VerifyResponse{
		Authorized: true,
		UserId:     user.UserID,
		RateLimit:  int32(user.RateLimit),
		Tier:       user.Tier,
	}, nil
}
