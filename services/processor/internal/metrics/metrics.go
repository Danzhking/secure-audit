package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EventsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "processor_events_processed_total",
		Help: "Total events processed",
	}, []string{"status"})

	AlertsTriggered = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "processor_alerts_triggered_total",
		Help: "Total alerts triggered by detection engine",
	}, []string{"rule"})

	ProcessingDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "processor_event_duration_seconds",
		Help:    "Time to process a single event",
		Buckets: prometheus.DefBuckets,
	})
)
