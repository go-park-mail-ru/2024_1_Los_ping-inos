package image

type (
	Image struct {
		UserId     int64  `json:"person_id"`
		Url        string `json:"image_url"`
		CellNumber string `json:"cell"`
		FileName   string `json:"filename"`
	}
)
