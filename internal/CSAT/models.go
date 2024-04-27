package CSAT

type (
	CreateRequest struct {
		Q1       int `json:"q1"`
		TittleID int `json:"tittle"`
	}

	StatRequest struct {
		Tittle int `json:"tittle"`
	}

	StatResponse struct {
		//Q1Stat map[string]int `json:"q1Stat"`
		Tittle  string  `json:"tittle"`
		Average float32 `json:"avg"`
		Stats   []int   `json:"stats"`
	}
)
