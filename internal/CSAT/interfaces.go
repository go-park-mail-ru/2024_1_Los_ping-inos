package CSAT

import (
	"context"
)

type (
	UseCase interface {
		Create(ctx context.Context, request CreateRequest) error
		GetStat(ctx context.Context) ([][]string, error)
	}

	CsatStorage interface {
		Create(ctx context.Context, q1 int, tittleID int) error
		GetStat(ctx context.Context, tittleID int) (string, float32, []int, error)
		GetTittlesCount(ctx context.Context) (int, error)
	}
)
