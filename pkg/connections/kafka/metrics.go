package kafka

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/max-messenger/max-bot-example-todolist/pkg/telemetry"
)

var (
	producerTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kafka_pool",
			Subsystem: "producer",
			Name:      "send_count_total",
			Help:      "Produced events count",
		},
		[]string{"name", "topic", telemetry.ErrLabel},
	)

	producerAsyncHandle = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "kafka_pool",
			Subsystem: "producer_async",
			Name:      "duration_seconds",
			Help:      "Time elapsed to async produce single messsage",
			Buckets:   telemetry.DefaultHistogramBuckets,
		},
		[]string{"name", "topic", telemetry.ErrLabel},
	)

	producerAsyncMessagesInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kafka_pool",
			Subsystem: "producer_async",
			Name:      "messages_in_flight",
			Help:      "Number of produced messages in flight",
		},
		[]string{"name", "topic"},
	)

	consumerHandle = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "kafka_pool",
			Subsystem: "consumer",
			Name:      "duration_seconds",
			Help:      "Time elapsed to consume single messsage",
			Buckets:   telemetry.DefaultHistogramBuckets,
		},
		[]string{"name", "topic", telemetry.ErrLabel},
	)

	consumerGroupHandle = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "kafka_pool",
			Subsystem: "consumer_group",
			Name:      "duration_seconds",
			Help:      "Time elapsed to consume single messsage",
			Buckets:   telemetry.DefaultHistogramBuckets,
		},
		[]string{"name", "group", "topic", telemetry.ErrLabel},
	)

	consumerGroupBatchHandle = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "kafka_pool",
			Subsystem: "consumer_group",
			Name:      "batch_duration_seconds",
			Help:      "Time elapsed to consume batch messsages",
			Buckets:   telemetry.DefaultHistogramBuckets,
		},
		[]string{"name", "group", "topic", telemetry.ErrLabel},
	)

	consumerGroupPollTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kafka",
			Subsystem: "consumer_group",
			Name:      "poll_total",
			Help:      "Amount of poll received per partition",
		},
		[]string{"name", "group", "topic", "partition", "is_full"},
	)
	consumerGroupBatchPollTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kafka",
			Subsystem: "consumer_group",
			Name:      "batch_poll_total",
			Help:      "Amount of poll received per partition",
		},
		[]string{"name", "group", "topic", "partition", "is_full"},
	)
)
