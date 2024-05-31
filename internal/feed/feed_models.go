package feed

import (
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"main.go/internal/types"
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
	InterestGetFilter struct {
		ID   []types.InterestID
		Name []string
	}
	Interest struct {
		ID   types.InterestID
		Name string `json:"name"`
	}
	ImageToSend struct {
		Cell string `json:"cell"`
		Url  string `json:"url"`
	}
	LikeGetFilter struct {
		Person1 *types.UserID
	}
	Image struct {
		UserId     int64  `json:"person_id"`
		Url        string `json:"image_url"`
		CellNumber string `json:"cell"`
	}
	Like struct {
		Person1 types.UserID
		Person2 types.UserID
	}

	//easyjson:json
	MessageToReceive struct {
		MsgType    string `json:"Type"`
		Properties struct {
			Id       int64        `json:"id"`
			Data     string       `json:"data"`
			Sender   types.UserID `json:"sender"`
			Receiver types.UserID `json:"receiver"`
			Time     int64        `json:"time"`
		} `json:"Properties"`
	}

	//easyjson:json
	Message struct {
		MsgType    string `json:"Type"`
		Properties MsgProperties
	}

	//easyjson:json
	MsgProperties struct {
		Id       int64        `json:"id"`
		Data     string       `json:"data"`
		Sender   types.UserID `json:"sender"`
		Receiver types.UserID `json:"receiver"`
		Time     time.Time    `json:"time"`
	}

	//easyjson:json
	GetChatRequest struct {
		Person types.UserID `json:"person"`
	}

	ChatPreview struct {
		PersonID    int64   `json:"personID"`
		Name        string  `json:"name"`
		Photo       string  `json:"photo"`
		LastMessage Message `json:"lastMessage"`
	}

	//easyjson:json
	CreateLikeRequest struct {
		Profile2 types.UserID `json:"id"`
	}

	//easyjson:json
	AllChats struct {
		Chats []ChatPreview `json:"chats"`
	}

	Claim struct {
		Id         int64 `json:"id"`
		TypeID     int64 `json:"typeID"`
		SenderID   int64 `json:"senderID"`
		ReceiverID int64 `json:"receiverID"`
	}

	//easyjson:json
	CreateClaimRequest struct {
		Type       int64 `json:"type"`
		ReceiverID int64 `json:"receiverID"`
	}

	PureClaim struct {
		Id    int64  `json:"id"`
		Title string `json:"title"`
	}

	//easyjson:json
	CardsToSend struct {
		Cards []Card `json:"cards"`
	}

	//easyjson:json
	MessagesToSend struct {
		Messages []Message `json:"messages"`
	}

	//easyjson:json
	GetChatFull struct {
		Messages []Message `json:"messages"`
		Images   []Image   `json:"images"`
		Person   Person    `json:"person"`
	}

	//easyjson:json
	ClaimsToSend struct {
		Claims []PureClaim `json:"claims"`
	}
)

var (
	NoLikesLeftErr  = errors.New("no likes left")
	NoMatchFoundErr = errors.New("sql: Rows are closed")
	WSClosedErr     = errors.New("websocket: close sent")
	TotalHits       = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "feed_total_hits",
			Help: "Count of hits in feed service.",
		},
		[]string{},
	)
	HitDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "feed_methods_handling_duration",
			Help: "Duration processing hit",
		},
		[]string{"method", "path"},
	)
)
