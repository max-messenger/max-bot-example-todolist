package analytic

import (
	"go.uber.org/fx"

	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/kafka"
)

var Module = fx.Options(
	fx.Provide(
		New,
		producerAdapter,
	),
)

func producerAdapter(
	cfg Config,
	kPool *kafka.Pool,
) (Producer, error) {
	if !cfg.Enabled {
		return nil, nil // nolint:nilnil
	}

	kf, err := kPool.GetPool("analytic")
	if err != nil {
		return nil, err
	}

	return kafka.NewProducer(kf)
}
