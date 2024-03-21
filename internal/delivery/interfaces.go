package delivery

type Service interface {
	GetCards(sessionID string) (string, error)
	GetName(sessionID string) (string, error)
	GetAllInterests() (string, error)
}

type Auth interface {
	IsAuthenticated(sessionID string) bool
	Login(email, password string) (string, string, error)
	Logout(sessionID string) error
	Registration(Name string, Birthday string, Gender string, Email string, Password string) (string, string, error)
}
