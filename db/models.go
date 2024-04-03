package models

import (
	"main.go/internal/types"
	"time"
)

// Person model info
// @Description Информация об аккаунте пользователя
type Person struct {
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

type PersonGetFilter struct {
	ID        []types.UserID
	Birthday  []time.Time
	Location  []string
	Email     []string
	Premium   []bool
	SessionID []string
}

type InterestGetFilter struct {
	ID   []types.InterestID
	Name []string
}

type Interest struct {
	ID   types.InterestID
	Name string
}

type PersonInterest struct {
	PersonId   types.UserID
	InterestID types.InterestID
}

type PersonWithInterests struct {
	Person    *Person
	Interests []*Interest
}
