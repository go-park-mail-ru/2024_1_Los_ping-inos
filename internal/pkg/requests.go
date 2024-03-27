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
		Name        string   `json:"name"`
		Birthday    string   `json:"birthday"`
		Password    string   `json:"password"`
		Description string   `json:"description"`
		Interests   []string `json:"interests"`
	}
)
