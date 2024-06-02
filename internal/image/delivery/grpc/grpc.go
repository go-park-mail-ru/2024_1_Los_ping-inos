package delivery

import (
	"context"
	"main.go/internal/image"
	pb "main.go/internal/image/protos/gen"
	"main.go/internal/image/usecase"
	. "main.go/internal/logs"
	"time"
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
	start := time.Now()
	tmp := server.ctx.Value(Logg).(Log)
	tmp.RequestID = req.RequestID
	server.ctx = context.WithValue(server.ctx, Logg, tmp)
	println("THIS IS ID AND CELL", req.Id, req.Cell)
	url, err := server.useCase.GetImage(req.Id, req.Cell, server.ctx)

	end := time.Since(start)
	image.TotalHits.WithLabelValues().Inc()
	image.HitDuration.WithLabelValues("grpc", "getImage").Set(float64(end.Milliseconds()))
	return &pb.GetImageResponce{Url: url}, err
}
