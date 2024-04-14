package auth

import (
	"context"
	"main.go/internal/types"
)

type (
	IUseCase interface {
		IsAuthenticated(sessionID string, ctx context.Context) (types.UserID, bool)
		Login(email, password string, ctx context.Context) (string, types.UserID, error)
		Logout(sessionID string, ctx context.Context) error
		Registration(body RegitstrationBody, ctx context.Context) (string, types.UserID, error)
		GetAllInterests(ctx context.Context) ([]*Interest, error)
		GetName(sessionID string, ctx context.Context) (string, error)
	}
	PostgresRepo interface {
		AddAccount(ctx context.Context, Name string, Birhday string, Gender string, Email string, Password string) error
		Get(ctx context.Context, filter *PersonGetFilter) ([]*Person, error)
		Update(ctx context.Context, person Person) error
		RemoveSession(ctx context.Context, sid string) error
	}

	InterestStorage interface {
		Get(ctx context.Context, filter *InterestGetFilter) ([]*Interest, error)
		CreatePersonInterests(ctx context.Context, personID types.UserID, interestID []types.InterestID) error
	}
)
