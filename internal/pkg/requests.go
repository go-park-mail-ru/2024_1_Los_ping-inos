package requests

type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RegistrationRequest struct {
		Name     string `json:"name"`
		Birthday string `json:"birthday"`
		Gender   string `json:"gender"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)
