package delivery

import (
	"context"

	pb "main.go/internal/image/protos/gen"
	"main.go/internal/image/usecase"
	. "main.go/internal/logs"
)

type Server struct {
	pb.UnimplementedImageServer
	useCase *usecase.UseCase
	ctx     context.Context
}

func NewGRPCDeliver(uc *usecase.UseCase) *Server {
	res := &Server{useCase: uc}
	logger := InitLog()
	res.ctx = context.WithValue(context.Background(), Logg, logger)
	return res
}
