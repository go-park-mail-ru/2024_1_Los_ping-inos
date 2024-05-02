package delivery

import (
	"context"

	pb "main.go/internal/auth/proto"
	"main.go/internal/auth/usecase"
	. "main.go/internal/logs"
)

type Server struct {
	pb.UnimplementedAuthHandlServer
	useCase *usecase.UseCase
	ctx     context.Context
}

func NewGRPCDeliver(uc *usecase.UseCase) *Server {
	res := &Server{useCase: uc}
	logger := InitLog()
	res.ctx = context.WithValue(context.Background(), Logg, logger)
	return res
}

//func NewServer()

func (server *Server) IsAuthenticated(_ context.Context, req *pb.IsAuthRequest) (*pb.IsAuthResponse, error) {
	tmp := server.ctx.Value(Logg).(Log) // TODO intercepter
	tmp.RequestID = req.RequestID
	server.ctx = context.WithValue(server.ctx, Logg, tmp)
	res, ok, err := server.useCase.IsAuthenticated(req.SessionID, server.ctx)
	return &pb.IsAuthResponse{IsAuthenticated: ok, UserID: int64(res)}, err
}
