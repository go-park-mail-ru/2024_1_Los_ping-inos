package CSAT

import (
	"context"
)

type (
	UseCase interface {
		Create(ctx context.Context, request CreateRequest) error
		GetStat(ctx context.Context) (map[string]int, error)
	}

	CsatStorage interface {
		Create(ctx context.Context, q1 int) error
		GetStat(ctx context.Context) (map[string]int, error)
	}
)
