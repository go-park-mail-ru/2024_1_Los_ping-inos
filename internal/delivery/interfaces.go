package delivery

import "main.go/internal/types"

type Service interface {
	GetCards(sessionID string, firstID types.UserID) (string, error)
	GetAllInterests() (string, error)
}

type Auth interface {
	IsAuthenticated(sessionID string) bool
	Login(email, password string) (string, error)
	Logout(sessionID string) error
	Registration(Name string, Birthday string, Gender string, Email string, Password string) (string, error)
}
