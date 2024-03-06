package models

import (
	"main.go/internal/types"
	"time"
)

type Person struct {
	ID          types.UserID `json:"ID"`
	Name        string       `json:"name"`
	Birthday    time.Time    `json:"birthday"`
	Description string       `json:"description"`
	Location    string       `json:"-"`
	Photo       string       `json:"photo"`
	Email       string       `json:"-"`
	Password    string       `json:"-"`
	Gender      string       `json:"gender"`
	CreatedAt   time.Time    `json:"-"`
	Premium     bool         `json:"-"`
	LikesLeft   int          `json:"-"`
	SessionID   string       `json:"session_id"`
}

type PersonGetFilter struct {
	ID        []types.UserID
	Birthday  []time.Time
	Location  []string
	Email     []string
	Premium   []bool
	SessionID []string
}

type Interest struct {
	ID   types.InterestID
	Name string
}
