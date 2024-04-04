package service

import (
	"main.go/internal/types"
)

type ProfileUpdate struct {
	SessionID   string
	Name        string
	Email       string
	Password    string
	Description string
	Birthday    string
	Interests   []string
}

type ProfileGetParams struct {
	ID        []types.UserID
	SessionID []string
}
