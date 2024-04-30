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

	ProfileGetParams struct {
		ID        []types.UserID
		SessionID []string
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

	Profile struct {
		ID          types.UserID  `json:"-"`
		Name        string        `json:"name"`
		Birthday    time.Time     `json:"birthday"`
		Description string        `json:"description"`
		Email       string        `json:"email"`
		Interests   []*Interest   `json:"interests"`
		Photos      []ImageToSend `json:"photos"`
		CSRFT       string        `json:"csrft"`
	}

	Image struct {
		UserId     int64  `json:"person_id"`
		Url        string `json:"image_url"`
		CellNumber string `json:"cell"`
	}

	ImageToSend struct {
		Cell string `json:"cell"`
		Url  string `json:"url"`
	}

	ProfileUpdateRequest struct {
		Name        string   `json:"name"`
		Email       string   `json:"email"`
		Birthday    string   `json:"birthday"`
		Password    string   `json:"password"`
		OldPassword string   `json:"oldPassword"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
	}

	Session struct {
		UID types.UserID `json:"UID"`
		SID string       `json:"SID"`
	}
)
