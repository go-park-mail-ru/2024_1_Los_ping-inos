package models

import (
	"main.go/internal/types"
	"time"
)

type Person struct {
	ID          types.UserID
	Name        string
	Birthday    time.Time
	Description string
	Location    string
	Photo       string
	Email       string
	Password    string
	CreatedAt   time.Time
	Premium     bool
	LikesLeft   int
	SessionID   string
}

type PersonFilter struct {
	ID        []types.UserID
	Birthday  []time.Time
	Location  []string
	Email     []string
	Premium   []bool
	SessionID []string
}
