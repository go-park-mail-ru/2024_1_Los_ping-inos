package auth

import (
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"main.go/internal/types"
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
		ID             types.UserID `json:"ID"`
		Name           string       `json:"name"`
		Birthday       time.Time    `json:"birthday"`
		Description    string       `json:"description"`
		Location       string       `json:"-"`
		Photo          string       `json:"photo"`
		Email          string       `json:"email"`
		Password       string       `json:"-"`
		Gender         string       `json:"gender"`
		CreatedAt      time.Time    `json:"-"`
		Premium        bool         `json:"premium"`
		PremiumExpires int64        `json:"premiumExpires"`
		LikesLeft      int          `json:"-"`
		SessionID      string       `json:"session_id"`
	}

	ProfileGetParams struct {
		ID        []types.UserID
		SessionID []string
		Name      string
		NeedEmail bool
	}

	PersonGetFilter struct {
		ID        []types.UserID
		Email     []string
		SessionID []string
		Name      string
	}

	//easyjson:json
	Interest struct {
		ID   types.InterestID
		Name string `json:"name"`
	}

	InterestGetFilter struct {
		ID   []types.InterestID
		Name []string
	}

	//easyjson:json
	Profile struct {
		ID             types.UserID  `json:"id"`
		Name           string        `json:"name"`
		Birthday       time.Time     `json:"birthday"`
		Description    string        `json:"description"`
		Email          string        `json:"email"`
		Premium        bool          `json:"premium"`
		PremiumExpires int64         `json:"premiumExpires"`
		LikesLeft      int           `json:"likesLeft"`
		Interests      []*Interest   `json:"interests"`
		Photos         []ImageToSend `json:"photos"`
		CSRFT          string        `json:"csrft"`
	}

	//easyjson:json
	CSRFTokenResponse struct {
		Csrft string `json:"csrft"`
	}

	//easyjson:json
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//easyjson:json
	RegistrationRequest struct {
		Name      string   `json:"name"`
		Birthday  string   `json:"birthday"`
		Gender    string   `json:"gender"`
		Email     string   `json:"email"`
		Password  string   `json:"password"`
		Interests []string `json:"interests"`
	}

	Image struct {
		UserId     int64  `json:"person_id"`
		Url        string `json:"image_url"`
		CellNumber string `json:"cell"`
	}

	//easyjson:json
	ImageToSend struct {
		Cell string `json:"cell"`
		Url  string `json:"url"`
	}

	//easyjson:json
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

	//easyjson:json
	GetMatchesRequest struct {
		Name string `json:"name"`
	}

	//easyjson:json
	Operations struct {
		Operations []PaymentOperation `json:"operations"`
	}

	//easyjson:json
	PaymentOperation struct {
		Label    string    `json:"label"`
		Datetime time.Time `json:"datetime"`
	}

	//easyjson:json
	PaymentHistory struct {
		Times []HistoryRecord `json:"records"`
	}

	HistoryRecord struct {
		Time  int64  `json:"time"`
		Sum   string `json:"sum"`
		Title string `json:"title"`
	}

	//easyjson:json
	Matches struct {
		Match []Profile `json:"matches"`
	}
	//easyjson:json
	Interests struct {
		Interes []*Interest `json:"interests"`
	}
)

var (
	NoPaymentErr = errors.New("no payment got")

	TotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_total_hits",
			Help: "Count of hits in auth service.",
		},
		[]string{},
	)
	HitDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "auth_methods_handling_duration",
			Help: "Duration processing hit",
		},
		[]string{"method", "path"},
	)
)
