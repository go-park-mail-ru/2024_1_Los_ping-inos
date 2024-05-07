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

func (server *Server) GetImage(_ context.Context, req *pb.GetImageRequest) (*pb.GetImageResponce, error) {
	tmp := server.ctx.Value(Logg).(Log)
	tmp.RequestID = req.RequestID
	server.ctx = context.WithValue(server.ctx, Logg, tmp)
	println("THIS IS ID AND CELL", req.Id, req.Cell)
	url, err := server.useCase.GetImage(req.Id, req.Cell, server.ctx)
	println("THIS IS ULR", url)
	return &pb.GetImageResponce{Url: url}, err
}
