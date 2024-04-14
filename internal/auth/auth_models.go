package auth

import (
	"main.go/internal/types"
	"time"
)

type (
	RegitstrationBody struct {
		Name      string
		Birthday  string
		Gender    string
		Email     string
		Password  string
		Interests []string
	}

	Person struct {
		ID          types.UserID `json:"ID"`
		Name        string       `json:"name"`
		Birthday    time.Time    `json:"birthday"`
		Description string       `json:"description"`
		Location    string       `json:"-"`
		Photo       string       `json:"photo"`
		Email       string       `json:"email"`
		Password    string       `json:"-"`
		Gender      string       `json:"gender"`
		CreatedAt   time.Time    `json:"-"`
		Premium     bool         `json:"-"`
		LikesLeft   int          `json:"-"`
		SessionID   string       `json:"session_id"`
	}

	PersonGetFilter struct {
		ID        []types.UserID
		Email     []string
		SessionID []string
	}

	Interest struct {
		ID   types.InterestID
		Name string `json:"name"`
	}

	InterestGetFilter struct {
		ID   []types.InterestID
		Name []string
	}
)
