package usecase

import (
	"context"
	"fmt"

	"main.go/internal/CSAT"
)

type UseCase struct {
	stor CSAT.CsatStorage
}

func NewUseCase(cstore CSAT.CsatStorage) *UseCase {
	return &UseCase{
		stor: cstore,
	}
}

func (service *UseCase) Create(ctx context.Context, request CSAT.CreateRequest) error {
	return service.stor.Create(ctx, request.Q1, request.TittleID)
}

func (service *UseCase) GetStat(ctx context.Context) ([][]string, error) {
	//tittle, avg, stats, err := service.stor.GetStat(ctx, tittleID)
	count, err := service.stor.GetTittlesCount(ctx)
	if err != nil {
		return nil, err
	}

	var allQ [][]string

	for count > 0 {
		tittle, avg, stats, err := service.stor.GetStat(ctx, count)
		if err != nil {
			return nil, err
		}
		var resp []string

		resp = append(resp, tittle)
		resp = append(resp, fmt.Sprintf("%f", avg))

		for _, stat := range stats {
			//println(stat)
			resp = append(resp, fmt.Sprintf("%v", stat))
		}

		allQ = append(allQ, resp)
		count--
	}

	return allQ, nil
}
