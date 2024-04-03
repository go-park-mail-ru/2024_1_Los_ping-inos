package service

type ProfileUpdate struct {
	SessionID   string
	Name        string
	Email       string
	Password    string
	Description string
	Birthday    string
	Interests   []string
}
