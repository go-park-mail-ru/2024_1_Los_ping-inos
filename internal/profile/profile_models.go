package profile

import (
	"main.go/internal/types"
	"time"
)

type (
	ProfileGetParams struct {
		ID        []types.UserID
		SessionID []string
	}
	// Card model info
	// @Description Информация в профиле пользователя (данные пользователя и его интересы)
	// имя возраст описание интересы фотографии
	Card struct {
		ID          types.UserID  `json:"id"`
		Name        string        `json:"name"`
		Birthday    time.Time     `json:"birthday"`
		Description string        `json:"description"`
		Email       string        `json:"email"`
		Interests   []*Interest   `json:"interests"`
		Photos      []ImageToSend `json:"photos"`
	}
	Interest struct {
		ID   types.InterestID
		Name string `json:"name"`
	}
	ImageToSend struct {
		Cell string `json:"cell"`
		Url  string `json:"url"`
	}
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
	Image struct {
		UserId     int64  `json:"person_id"`
		Url        string `json:"image_url"`
		CellNumber string `json:"cell"`
	}
	LikeGetFilter struct {
		Person1 *types.UserID
	}
	Like struct {
		Person1 types.UserID
		Person2 types.UserID
	}
)
