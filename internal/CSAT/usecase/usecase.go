package usecase

import (
	"context"
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
	return service.stor.Create(ctx, request.Q1)
}

func (service *UseCase) GetStat(ctx context.Context) (map[string]int, error) {
	return service.stor.GetStat(ctx)
}
