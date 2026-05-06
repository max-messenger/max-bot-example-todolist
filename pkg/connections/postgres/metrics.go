package postgres

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	duration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "postgres_request_duration_seconds",
			Help:    "postgres request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"action", "addr", "error"},
	)
	connections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "postgres_pool_connections",
			Help: "postgres pool connections",
		},
		[]string{"addr", "status"},
	)
)
