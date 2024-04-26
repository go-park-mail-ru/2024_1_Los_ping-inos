package delivery

import (
	"context"
	pb "main.go/internal/auth/proto"
	"main.go/internal/auth/usecase"
)

type Server struct {
	pb.UnimplementedAuthHandlServer
	useCase *usecase.UseCase
}

func NewGRPCDeliver(uc *usecase.UseCase) *Server {
	return &Server{
		useCase: uc,
	}
}

func (server *Server) CheckIsAuth(ctx context.Context, req *pb.IsAuthRequest) (*pb.IsAuthResponse, error) {
	res, ok, err := server.useCase.IsAuthenticated(req.SessionID, ctx)
	return &pb.IsAuthResponse{IsAuthenticated: ok, UserID: int64(res)}, err
}
