package delivery

type (
	Service interface {
		GetCards(sessionID string, reqID int64) (string, error)
		GetName(sessionID string, reqID int64) (string, error)
		GetAllInterests(reqID int64) (string, error)
		GetProfile(sessionID string, requestID int64) (string, error)
		UpdateProfile(sessionID, name, password, description, birthday string, interests []string, requestID int64) error
	}

	Auth interface {
		IsAuthenticated(sessionID string, reqID int64) bool
		Login(email, password string, reqID int64) (string, error)
		Logout(sessionID string, reqID int64) error
		Registration(Name string, Birthday string, Gender string, Email string, Password string, reqID int64) (string, error)
	}
)
