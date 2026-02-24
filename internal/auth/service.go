package auth

import (
    "context"
    "github.com/mohammad-kh1/distributed-gateway/api/proto"
)

type Server struct {
    proto.UnimplementedAuthServiceServer
}

func (s *Server) VerifyToken(ctx context.Context, req *proto.VerifyRequest) (*proto.VerifyResponse, error) {
    // TODO: real token validation (JWT, DB lookup, etc.)
    // For now — fake/test data

    tokens := map[string]*proto.VerifyResponse{  // ← note: *VerifyResponse
        "user-token-1": {
            Authorized: true,
            UserId:     "user_101",
            RateLimit:  5,
            Tier:       "free",
        },
        "admin-token": {
            Authorized: true,
            UserId:     "admin_001",
            RateLimit:  100,
            Tier:       "premium",
        },
    }

    res, ok := tokens[req.Token]
    if !ok {
        return &proto.VerifyResponse{
            Authorized: false,
            // optionally: UserId: "", RateLimit: 0, Tier: "none"
        }, nil
    }

    return res, nil
}
