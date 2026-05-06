package grpccli

import (
	"github.com/prometheus/client_golang/prometheus"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"

	"github.com/max-messenger/max-bot-example-todolist/pkg/telemetry"
)

var (
	clMetrics = grpcprom.NewClientMetrics(
		grpcprom.WithClientHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets(telemetry.DefaultHistogramBuckets),
		),
	)
)

func init() {
	prometheus.MustRegister(clMetrics)
}
