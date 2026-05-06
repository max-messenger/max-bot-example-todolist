# BGTasker

Асинхронное выполнение задач в фоне с использованием BGTasker.

## Концепция

BGTasker - это механизм для выполнения задач асинхронно в фоновых горутинах с поддержкой graceful shutdown.

Поддерживает два режима выполнения:

- **Drop** - отбрасывает задачи если очередь заполнена
- **Block** - блокирует до освобождения места в очереди

## Конфигурация

```yaml
bgtasker:
  tasks:
    # Конфигурация для аналитики
    analytic:
      runner_type: "drop"          # или "block"
      async_workers_num: 10         # количество воркеров
      async_buffer_size: 1000       # размер буфера очереди
      async_worker_timeout: 2s       # timeout выполнения задачи

    # Конфигурация для обработки событий
    events:
      runner_type: "block"
      async_workers_num: 20
      async_buffer_size: 5000
      async_worker_timeout: 5s

    # Дополнительные tasker'ы
    notifications:
      runner_type: "drop"
      async_workers_num: 5
      async_buffer_size: 500
      async_worker_timeout: 1s
```

## Использование BGTasker

### Получение Runner

```go
package analytics

import (
    "context"
    "github.com/max-messenger/max-bot-example-todolist/pkg/bgtasker"
)

type Service struct {
    tasker *bgtasker.Runner
}

func New(pool *bgtasker.Pool) *Service {
    // Получаем tasker с именем "analytic"
    tasker, err := pool.Get("analytic")
    if err != nil {
        panic(err)
    }

    return &Service{
        tasker: tasker,
    }
}
```

### Выполнение задач

```go
// Асинхронное выполнение (неблокирующее)
func (s *Service) SendEvent(ctx context.Context, event Event) {
    s.tasker.Run(ctx, func(taskCtx context.Context) error {
        // Задача будет выполнена в фоновой горутине
        return s.processEvent(taskCtx, event)
    })
}
```

### Обработка ошибок

```go
func (s *Service) SendEvent(ctx context.Context, event Event) {
    s.tasker.Run(ctx, func(taskCtx context.Context) error {
        if err := s.processEvent(taskCtx, event); err != nil {
            // Ошибка будет залогирована автоматически
            // но не вернется вызывающему коду
            return fmt.Errorf("process event: %w", err)
        }

        return nil
    })
}
```

## Режимы выполнения

### Drop Mode (по умолчанию)

Если очередь заполнена, задача отбрасывается с логированием ошибки:

```go
// Конфигурация
analytic:
  runner_type: "drop"
  async_buffer_size: 1000
```

```go
func (s *Service) SendEvent(ctx context.Context, event Event) {
    s.tasker.Run(ctx, func(taskCtx context.Context) error {
        // Если буфер заполнен, задача будет отброшена
        // и залогировано: "dropping event, channel is full"
        return s.processEvent(taskCtx, event)
    })
}
```

**Когда использовать:**

- События не критичны
- Можно потерять часть данных
- Нужно избежать блокировки вызывающего кода

### Block Mode

Если очередь заполнена, ожидает освобождения места:

```go
// Конфигурация
events:
  runner_type: "block"
  async_buffer_size: 5000
```

```go
func (s *Service) ProcessEvent(ctx context.Context, event Event) {
    s.tasker.Run(ctx, func(taskCtx context.Context) error {
        // Если буфер заполнен, будет блокироваться
        // до освобождения места
        return s.processEvent(taskCtx, event)
    })
}
```

**Когда использовать:**

- Все события критичны
- Нельзя потерять данные
- Допустима блокировка при высокой нагрузке

## Timeout

Каждая задача имеет timeout для выполнения:

```go
analytic:
  async_worker_timeout: 2s
```

```go
func (s *Service) SendEvent(ctx context.Context, event Event) {
    s.tasker.Run(ctx, func(taskCtx context.Context) error {
        // taskCtx будет отменен через 2 секунды
        // независимо от контекста вызывающего кода
        return s.processEvent(taskCtx, event)
    })
}
```

Обработка timeout:

```go
func (s *Service) processEvent(ctx context.Context, event Event) error {
    // Проверяем на timeout
    if ctx.Err() != nil {
        return ctx.Err()
    }

    // Долгая операция
    if err := s.someOperation(ctx, event); err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }

    return nil
}
```

## Полный пример с аналитикой

```go
package analytics

import (
    "context"
    "fmt"
    "time"

    "go.uber.org/zap"
    "github.com/max-messenger/max-bot-example-todolist/pkg/bgtasker"
)

type Event struct {
    UserID int64
    Action string
    Data   map[string]any
}

type Service struct {
    logger *zap.Logger
    tasker *bgtasker.Runner
}

func New(
    logger *zap.Logger,
    pool *bgtasker.Pool,
) *Service {
    tasker, err := pool.Get("analytic")
    if err != nil {
        panic(err)
    }

    return &Service{
        logger: logger.Named("analytics"),
        tasker: tasker,
    }
}

// SendEvent - неблокирующая отправка события
func (s *Service) SendEvent(ctx context.Context, userID int64, action string, data map[string]any) {
    event := Event{
        UserID: userID,
        Action: action,
        Data:   data,
    }

    s.tasker.Run(ctx, s.createTask(event))
}

// createTask - создает задачу для обработки события
func (s *Service) createTask(event Event) bgtasker.TaskFunc {
    return func(ctx context.Context) error {
        // Логируем начало обработки
        s.logger.Debug("processing analytic event",
            zap.Int64("user_id", event.UserID),
            zap.String("action", event.Action),
        )

        // Проверяем на timeout
        if ctx.Err() != nil {
            return ctx.Err()
        }

        // Обработка события
        if err := s.processEvent(ctx, event); err != nil {
            return fmt.Errorf("process event: %w", err)
        }

        // Успешное завершение
        s.logger.Debug("analytic event processed",
            zap.Int64("user_id", event.UserID),
            zap.String("action", event.Action),
        )

        return nil
    }
}

// processEvent - реальная логика обработки
func (s *Service) processEvent(ctx context.Context, event Event) error {
    // 1. Валидация
    if err := s.validateEvent(event); err != nil {
        return err
    }

    // 2. Обогащение данных
    enrichedData := s.enrichData(ctx, event)

    // 3. Отправка в Kafka/другую систему
    if err := s.sendToDestination(ctx, enrichedData); err != nil {
        return err
    }

    return nil
}

func (s *Service) validateEvent(event Event) error {
    if event.UserID <= 0 {
        return fmt.Errorf("invalid user_id: %d", event.UserID)
    }

    if event.Action == "" {
        return fmt.Errorf("empty action")
    }

    return nil
}

func (s *Service) enrichData(ctx context.Context, event Event) map[string]any {
    enriched := make(map[string]any, len(event.Data)+2)

    for k, v := range event.Data {
        enriched[k] = v
    }

    enriched["user_id"] = event.UserID
    enriched["action"] = event.Action
    enriched["timestamp"] = time.Now().Unix()

    return enriched
}

func (s *Service) sendToDestination(ctx context.Context, data map[string]any) error {
    // Реальная отправка в Kafka, analytics service и т.д.
    return nil
}
```

## FX модуль

```go
package analytics

import (
    "go.uber.org/fx"
    "go.uber.org/zap"
    "github.com/max-messenger/max-bot-example-todolist/pkg/bgtasker"
)

var Module = fx.Module(
    "analytics",
    fx.Provide(
        New,
    ),
    fx.Decorate(func(log *zap.Logger) *zap.Logger {
        return log.Named("analytics")
    }),
)
```

## Особенности BGTasker

### 1. Context Without Cancel

Внутри задачи используется `context.WithoutCancel`:

```go
func (r *Runner) dropRunner(ctx context.Context, taskFunc TaskFunc) {
    t := task{
        Context:  context.WithoutCancel(ctx), // отменяется только по timeout
        TaskFunc: taskFunc,
    }
    // ...
}
```

Это означает:

- Отмена внешнего контекста не остановит уже запущенную задачу
- Задача остановится только по timeout или завершению
- Позволяет выполнять длинные операции без прерывания

## Рекомендации

1. **Drop vs Block**:
   - Используйте **Drop** для некритичных событий (аналитика, метрики)
   - Используйте **Block** для критичных операций (обработка платежей)

2. **Timeout**:
   - Устанавливайте timeout чуть больше ожидаемого времени выполнения
   - Обрабатывайте `ctx.Err()` в долгих операциях

3. **Воркеры**:
   - Количество воркеров должно быть оптимальным для вашей нагрузки
   - Слишком много воркеров - высокий расход CPU
   - Слишком мало - задачи будут ждать в очереди

4. **Размер буфера**:
   - Больший буфер позволяет справляться с пиковыми нагрузками
   - Но увеличивает потребление памяти

5. **Логирование**:
   - Логируйте важные события внутри задачи
   - Ошибки логируются автоматически

6. **Обработка ошибок**:
   - Всегда возвращайте ошибку из задачи для корректного логирования
   - Не паникуйте внутри задачи

## Типичные паттерны

### Отправка аналитики (неблокирующая)

```go
func (s *Service) TrackAction(ctx context.Context, userID int64, action string) {
    s.tasker.Run(ctx, func(taskCtx context.Context) error {
        s.producer.SendAsync(context.WithoutCancel(taskCtx), "analytics", kafka.Payload{
            Value: s.encodeAction(userID, action),
        })
        return nil
    })
}
```

### Обработка событий (неблокирующая)

```go
func (s *Service) ProcessUserEvent(ctx context.Context, event UserEvent) {
    s.tasker.Run(ctx, func(taskCtx context.Context) error {
        // Сложная обработка
        if err := s.validateEvent(taskCtx, event); err != nil {
            return err
        }

        if err := s.updateDatabase(taskCtx, event); err != nil {
            return err
        }

        if err := s.sendNotifications(taskCtx, event); err != nil {
            return err
        }

        return nil
    })
}
```

### Агрегация данных (неблокирующая)

```go
func (s *Service) AggregateStats(ctx context.Context) {
    s.tasker.Run(ctx, func(taskCtx context.Context) error {
        // Сбор статистики
        stats, err := s.collectStats(taskCtx)
        if err != nil {
            return err
        }

        // Агрегация
        aggregated := s.aggregate(stats)

        // Сохранение результата
        return s.saveAggregated(taskCtx, aggregated)
    })
}
```

## Отличия от Grace

| BGTasker | Grace |
| -------- | ----- |
| Выполнение задач в воркерах | Жизненный цикл сервисов |
| Поддержка очереди | Запуск/остановка горутин |
| Timeout для задач | Context для остановки |
| Drop/Block режимы | Последовательная остановка |
| Метрики по задачам | Метрики по сервисам |

**Используйте BGTasker**, когда:

- Нужно выполнять много однотипных задач параллельно
- Неблокирующая отправка операций
- Есть очередь задач

**Используйте Grace**, когда:

- Нужно управлять жизненным циклом сервиса
- Горутина должна работать постоянно
- Нужен graceful shutdown
