package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsController = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "marathon_app_http_requests",
			Help: "Total number of HTTP requests",
		},
	)

	GetRunnerHttpResponsesCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "marathon_app_get_runner_http_responses",
			Help: "Total number of HTTP responses for GET /runner API",
		},
		[]string{"status"},
	)

	GetAllRunnersTimer = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name: "marathon_app_get_all_runners_duration",
			Help: "Duration of get all runners operation",
		},
	)
)
