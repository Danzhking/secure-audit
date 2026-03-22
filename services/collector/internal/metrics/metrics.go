package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EventsReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "collector_events_received_total",
		Help: "Total events received by the Collector",
	}, []string{"status"})

	EventsPublished = promauto.NewCounter(prometheus.CounterOpts{
		Name: "collector_events_published_total",
		Help: "Total events published to RabbitMQ",
	})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "collector_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	RateLimitRejected = promauto.NewCounter(prometheus.CounterOpts{
		Name: "collector_rate_limit_rejected_total",
		Help: "Requests rejected by rate limiter",
	})
)
