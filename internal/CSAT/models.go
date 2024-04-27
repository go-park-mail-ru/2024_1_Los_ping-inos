package CSAT

type (
	CreateRequest struct {
		Q1 int `json:"q1"`
	}

	StatResponse struct {
		Q1Stat map[string]int `json:"q1Stat"`
	}
)
