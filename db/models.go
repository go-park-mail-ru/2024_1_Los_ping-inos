package models

import (
	"main.go/internal/types"
	"time"
)

type (
	// Person model info
	// @Description Информация об аккаунте пользователя
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
		Birthday  []time.Time
		Location  []string
		Email     []string
		Premium   []bool
		SessionID []string
	}

	InterestGetFilter struct {
		ID   []types.InterestID
		Name []string
	}

	LikeGetFilter struct {
		Person1 *types.UserID
		Person2 *types.UserID
	}

	Like struct {
		Person1 types.UserID
		Person2 types.UserID
	}

	Interest struct {
		ID   types.InterestID
		Name string
	}

	PersonInterest struct {
		PersonId   types.UserID
		InterestID types.InterestID
	}

	// PersonWithInterests model info
	// @Description Информация в профиле пользователя (данные пользователя и его интересы)
	PersonWithInterests struct {
		Person    *Person     `json:"person"`
		Interests []*Interest `json:"interests"`
	}
)
