package delivery

import (
	"context"
	"github.com/sirupsen/logrus"
	pb "main.go/internal/auth/proto"
	"main.go/internal/auth/usecase"
	. "main.go/internal/logs"
	"main.go/internal/types"
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

func (server *Server) IsAuthenticated(_ context.Context, req *pb.IsAuthRequest) (*pb.IsAuthResponse, error) {
	tmp := server.ctx.Value(Logg).(Log) // TODO intercepter
	tmp.RequestID = req.RequestID
	server.ctx = context.WithValue(server.ctx, Logg, tmp)
	res, ok, err := server.useCase.IsAuthenticated(req.SessionID, server.ctx)
	return &pb.IsAuthResponse{IsAuthenticated: ok, UserID: int64(res)}, err
}

func (server *Server) GetMatches(_ context.Context, req *pb.GetMatchesRequest) (*pb.GetMatchesResponse, error) {
	tmp := server.ctx.Value(Logg).(Log) // TODO intercepter
	tmp.RequestID = req.RequestID
	server.ctx = context.WithValue(server.ctx, Logg, tmp)
	matches, err := server.useCase.GetMatches(types.UserID(req.UserID), server.ctx)
	if err != nil {
		tmp.Logger.WithFields(logrus.Fields{RequestID: tmp.RequestID}).Warn("can't get matches: ", err.Error())
		return nil, err
	}
	res := make([]*pb.Match, len(matches))
	for i := range matches {
		res[i] = &pb.Match{
			Name:     matches[i].Name,
			PersonID: int64(matches[i].ID),
		}

		if len(matches[i].Photos) > 0 {
			res[i].Photo = matches[i].Photos[0].Url
		}
	}
	return &pb.GetMatchesResponse{Matches: res}, nil
}
