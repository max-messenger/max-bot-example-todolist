package health

import (
	"sync/atomic"
)

type Health struct {
	cfg     Config
	isReady atomic.Bool // global readiness flag, initial is false
}

func NewHealth(cfg Config) *Health {
	return &Health{
		cfg: cfg,
	}
}

func (h *Health) SetReady(ready bool) {
	h.isReady.Store(ready)
}
