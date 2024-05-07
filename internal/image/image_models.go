package image

import "github.com/prometheus/client_golang/prometheus"

type (
	Image struct {
		UserId     int64  `json:"person_id"`
		Url        string `json:"image_url"`
		CellNumber string `json:"cell"`
		FileName   string `json:"filename"`
	}
)

var (
	TotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "images_total_hits",
			Help: "Count of hits in images service.",
		},
		[]string{},
	)
	HitDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "images_methods_handling_duration",
			Help: "Duration processing hit",
		},
		[]string{"method", "path"},
	)
)
