package requests

type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RegistrationRequest struct {
		Name      string   `json:"name"`
		Birthday  string   `json:"birthday"`
		Gender    string   `json:"gender"`
		Email     string   `json:"email"`
		Password  string   `json:"password"`
		Interests []string `json:"interests"`
	}

	ProfileUpdateRequest struct {
		SID         string   `json:"SID"`
		Name        string   `json:"name"`
		Email       string   `json:"email"`
		Birthday    string   `json:"birthday"`
		Password    string   `json:"password"`
		OldPassword string   `json:"oldPassword"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
	}
)
