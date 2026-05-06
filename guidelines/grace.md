# Работа с Grace

Механизм для управления жизненным циклом фоновых сервисов с поддержкой graceful shutdown.

## Что такое Grace?

`Grace` - это механизм управления сервисами, которые должны запускаться и останавливаться в процессе работы приложения. Все такие сервисы реализуют интерфейс:

```go
type Service interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}
```

## Создание фонового сервиса

### Базовый пример

На основе примера `Subscriber` из `internal/app/subscriber/subscriber.go`:

```go
package subscriber

import (
    "context"
    "go.uber.org/zap"
)

type Subscriber struct {
    logger *zap.Logger
    quit   chan struct{}
}

func New(logger *zap.Logger) *Subscriber {
    return &Subscriber{
        logger: logger,
        quit:   make(chan struct{}),
    }
}

// Start запускает фоновую горутину
func (s *Subscriber) Start(_ context.Context) error {
    go s.messagesConsumer()
    return nil
}

// Stop останавливает фоновую горутину
func (s *Subscriber) Stop(_ context.Context) error {
    close(s.quit)
    return nil
}

func (s *Subscriber) messagesConsumer() {
    for {
        select {
        case <-s.quit:
            return // выход из горутины
        case msg := <-s.messagesChannel:
            // обработка сообщений
            go func() {
                // обработка в отдельной горутине для неблокирующего чтения
                s.processMessage(msg)
            }()
        }
    }
}
```

## Регистрация сервиса в Grace

Для автоматического управления сервисом через Grace, необходимо адаптировать его в fx-модуле:

```go
package subscriber

import (
    "go.uber.org/fx"
    "github.com/max-messenger/max-bot-example-todolist/pkg/grace"
)

var Module = fx.Module(
    "subscriber",
    fx.Provide(
        New,
        graceAdapter, // адаптер для Grace
    ),
)
```

### Адаптер для Grace

```go
// graceOut - структура для экспорта сервиса в группу "grace"
type graceOut struct {
    fx.Out

    Service grace.Service `group:"grace"`
}

// graceAdapter адаптирует Subscriber к интерфейсу grace.Service
func graceAdapter(sub *Subscriber) graceOut {
    return graceOut{
        Service: sub,
    }
}
```

## Примеры использования

Периодическая задача

```go
package cleanup

import (
    "context"
    "time"
    "go.uber.org/zap"
)

type CleanupService struct {
    logger *zap.Logger
    ticker *time.Ticker
    quit   chan struct{}
}

func New(logger *zap.Logger) *CleanupService {
    return &CleanupService{
        logger: logger,
        ticker: time.NewTicker(5 * time.Minute),
        quit:   make(chan struct{}),
    }
}

func (cs *CleanupService) Start(_ context.Context) error {
    go cs.run()
    return nil
}

func (cs *CleanupService) Stop(_ context.Context) error {
    close(cs.quit)
    cs.ticker.Stop()
    return nil
}

func (cs *CleanupService) run() {
    for {
        select {
        case <-cs.quit:
            return
        case <-cs.ticker.C:
            if err := cs.cleanup(); err != nil {
                cs.logger.Error("cleanup failed", zap.Error(err))
            }
        }
    }
}
```

## Полный пример с FX модулем

```go
package myservice

import (
    "context"
    "go.uber.org/fx"
    "go.uber.org/zap"
    "github.com/max-messenger/max-bot-example-todolist/pkg/grace"
)

type MyService struct {
    logger *zap.Logger
    quit   chan struct{}
}

func New(logger *zap.Logger) *MyService {
    return &MyService{
        logger: logger.Named("myservice"),
        quit:   make(chan struct{}),
    }
}

func (s *MyService) Start(_ context.Context) error {
    go s.run()
    return nil
}

func (s *MyService) Stop(_ context.Context) error {
    close(s.quit)
    return nil
}

func (s *MyService) run() {
    for {
        select {
        case <-s.quit:
            s.logger.Info("stopping service")
            return
        case <-time.After(1 * time.Minute):
            s.logger.Info("service tick")
        }
    }
}

var Module = fx.Module(
    "myservice",
    fx.Provide(
        New,
        graceAdapter,
    ),
)

type graceOut struct {
    fx.Out
    Service grace.Service `group:"grace"`
}

func graceAdapter(svc *MyService) graceOut {
    return graceOut{Service: svc}
}
```
