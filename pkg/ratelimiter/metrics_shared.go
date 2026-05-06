package ratelimiter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/max-messenger/max-bot-example-todolist/pkg/telemetry"
)

var (
	sharedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shared_rate_limiter_actions_total",
			Help: "Total number of shared rate limiter action requests.",
		},
		[]string{"action", "status"},
	)

	sharedDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shared_rate_limiter_action_duration_seconds",
			Help:    "shared rate limiter action request latencies in seconds.",
			Buckets: telemetry.DefaultHistogramBuckets,
		},
		[]string{"action"},
	)

	sharedExceed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shared_rate_limiter_exceed_total",
			Help: "Total number of shared rate limiter action exceed.",
		},
		[]string{"action"},
	)

	sharedFails = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shared_rate_limiter_fails_total",
			Help: "Total number of shared rate limiter action fails.",
		},
		[]string{"action"},
	)
)
