# Чек-лист разработки нового сервиса

Пошаговый список действий при разработке нового сервиса в todolist.

## 1. Планирование

- [ ] Определить назначение сервиса
- [ ] Описать доменные сущности
- [ ] Определить требования к хранилищам (PostgreSQL, Redis, Kafka)
- [ ] Определить нужны ли фоновые задачи
- [ ] Составить список API endpoints

## 2. Domain слой

- [ ] Создать файлы доменных сущностей в `internal/app/domain/`
  - [ ] Определить структуры для сущностей

**Пример:**

```go
// internal/app/domain/todo.go
package domain

type Todo struct {
    ID        int64
    UserID    int64
    Message   string
    Done      bool
    CreatedAt time.Time
}
```

## 3. Service слой

- [ ] Создать папку сервиса в `internal/app/services/{service}/`
  - [ ] Определить интерфейс Repository (если треюбуется)
  - [ ] Создать структуру Service
  - [ ] Реализовать бизнес-методы

**Пример:**

```go
// internal/app/services/todolist/todolist.go
package todolist

//go:generate go tool mockgen -source=todolist.go -destination=todolist_mocks_test.go -package=todolist_test

type Repository interface {
    Add(ctx context.Context, d domain.Todo) error
    List(ctx context.Context) ([]domain.Todo, error)
}

type Service struct {
    repo Repository
}

func NewService(r Repository) *Service {
    return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, d domain.Todo) error {
    if err := d.Validate(); err != nil {
        return err
    }
    return s.repo.Add(ctx, d)
}
```

- [ ] Создать FX модуль `fx.go`

```go
// internal/app/services/todolist/fx.go
package todolist

import (
    "go.uber.org/fx"
    "go.uber.org/zap"
    "github.com/max-messenger/max-bot-example-todolist/internal/app/repository/todorepo"
)

var Module = fx.Module(
    "todolist",
    fx.Provide(
        NewService,
        todorepo.NewRepository,
    ),
    fx.Decorate(func(log *zap.Logger) *zap.Logger {
        return log.Named("todolist")
    }),
)
```

## 4. Repository слой

### Если используется PostgreSQL

- [ ] Создать папку репозитория в `internal/app/repository/{repo}/`
  - [ ] Определить константы для имен пулов
  - [ ] Создать структуру Repository с read/write полями
  - [ ] Реализовать методы Repository
  - [ ] Добавить entity структуру для маппинга из БД
  - [ ] Использовать `sqlbuilder` или любой другой для запросов, либо без

**Пример:**

```go
// internal/app/repository/todorepo/repository.go
package todorepo

const (
    ReadPGPool  = "sel"
    WritePGPool = "upd"
)

type Repository struct {
    read  *postgres.Postgres
    write *postgres.Postgres
}

func NewRepository(pool *postgres.Pool) (*Repository, error) {
    read, err := pool.GetPool(ReadPGPool)
    if err != nil {
        return nil, fmt.Errorf("get pool read: %w", err)
    }

    write, err := pool.GetPool(WritePGPool)
    if err != nil {
        return nil, fmt.Errorf("get pool write: %w", err)
    }

    return &Repository{read: read, write: write}, nil
}
```

- [ ] Создать FX модуль `fx.go`

```go
// internal/app/repository/todorepo/fx.go
package todorepo

import (
    "go.uber.org/fx"
    "github.com/max-messenger/max-bot-example-todolist/pkg/connections/postgres"
)

var Module = fx.Module(
    "todorepo",
    fx.Provide(
        NewRepository,
    ),
)
```

### Если используется Redis

- [ ] Создать крепозиторий
  - [ ] Определить методы
  - [ ] Использовать JSON для сложных структур
  - [ ] Устанавливать TTL для ключей

### Если используется Kafka Producer

- [ ] Создать сервис с kafka.Producer
- [ ] Реализовать методы отправки сообщений
- [ ] Использовать `SendAsync` для неблокирующей отправки

### Если используется Kafka Consumer

- [ ] Создать консьюмер с kafka.Consumer
- [ ] Реализовать обработку сообщений
- [ ] Если нужен Consumer Group - использовать kafka.ConsumerGroup
- [ ] Регистрировать в Grace (для фоновой обработки), лучше поместить в `subscriber`

## 5. Router слой (если нужен HTTP API)

- [ ] Создать папку контроллера в `internal/app/router/{ctrl}/`
  - [ ] Создать структуру Controller
  - [ ] Реализовать методы-обработчики (handlers)
  - [ ] Добавить регистрацию routes в методе Register

**Пример:**

```go
// internal/app/router/todolistctrl/todolistctrl.go
package todolistctrl

type Controller struct {
    logger  *zap.Logger
    service *todolist.Service
}

func NewController(
    logger *zap.Logger,
    service *todolist.Service,
) *Controller {
    return &Controller{
        logger:  logger,
        service: service,
    }
}

func (c *Controller) Register(r chi.Router) {
    r.Post("/todos", c.Create)
    r.Get("/todos", c.List)
}
```

- [ ] Создать FX модуль `fx.go`

```go
// internal/app/router/todolistctrl/fx.go
package todolistctrl

import (
    "go.uber.org/fx"
    "go.uber.org/zap"
    "github.com/max-messenger/max-bot-example-todolist/internal/app/router"
    "github.com/max-messenger/max-bot-example-todolist/internal/app/services/todolist"
)

var Module = fx.Module(
    "todolistctrl",
    fx.Provide(
        NewController,
        AdapterControllerOut,
    ),
)

type ControllerOut struct {
    fx.Out

    Controller router.Controller `group:"controller"`
}

func AdapterControllerOut(ctrl *Controller) ControllerOut {
    return ControllerOut{
        Controller: ctrl,
    }
}

```

## 6. Background Tasks (если нужны)

### Если нужен фоновый сервис (Grace)

- [ ] Создать сервис, реализующий `grace.Service`
  - [ ] Реализовать метод `Start` - запуск фоновой горутины
  - [ ] Реализовать метод `Stop` - graceful shutdown

**Пример:**

```go
type MyBackgroundService struct {
    logger *zap.Logger
    quit   chan struct{}
}

func (s *MyBackgroundService) Start(_ context.Context) error {
    go s.run()
    return nil
}

func (s *MyBackgroundService) Stop(_ context.Context) error {
    close(s.quit)
    return nil
}

func (s *MyBackgroundService) run() {
    for {
        select {
        case <-s.quit:
            return
        case <-time.After(1 * time.Minute):
            s.process()
        }
    }
}
```

- [ ] Создать адаптер для Grace

```go
type graceOut struct {
    fx.Out
    Service grace.Service `group:"grace"`
}

func graceAdapter(svc *MyBackgroundService) graceOut {
    return graceOut{Service: svc}
}
```

- [ ] Добавить адаптер в FX модуль

### Если нужны асинхронные задачи (BGTasker)

- [ ] Добавить конфигурацию tasker в `config.yaml`

```yaml
bgtasker:
  tasks:
    mytasks:
      runner_type: "drop"
      async_workers_num: 10
```

- [ ] Создать сервис с bgtasker.Runner
- [ ] Использовать `tasker.Run()` для выполнения задач
- [ ] Учитывать timeout при долгих операциях

**Пример:**

```go
type Service struct {
    tasker *bgtasker.Runner
}

func (s *Service) ProcessAsync(ctx context.Context, data Data) {
    s.tasker.Run(ctx, func(taskCtx context.Context) error {
        return s.processData(taskCtx, data)
    })
}
```

## 7. Конфигурация

- [ ] Добавить конфигурацию пулов PostgreSQL (если используется)

```yaml
postgres:
    sel:
        database_url: "postgres://..."
    upd:
        database_url: "postgres://..."
```

- [ ] Добавить конфигурацию пулов Redis (если используется)

```yaml
redis:
    default:
      addrs: ["localhost:6379"]
```

- [ ] Добавить конфигурацию Kafka (если используется)

```yaml
kafka:
  name:
    seeds: ["localhost:9092"]
```

- [ ] Добавить конфигурацию bgtasker (если используется)

```yaml
bgtasker:
  tasks:
    mytasks:
      runner_type: "drop"
```

## 8. Регистрация модуля

- [ ] Зарегистрировать модуль в `internal/app/fx.go`

```go
var Modules = fx.Options(
    // ...
    services.Module,      // включает ваш сервис
    repository.Module,    // включает ваш репозиторий
    // ...
)
```

## Частые ошибки

1. **Забыть зарегистрировать модуль** - сервис не будет создан
2. **Не использовать read/write разделение** - нагрузка на primary
3. **Не реализовать Grace интерфейс** - сервис не остановится корректно
4. **Не установить TTL в Redis** - утечка памяти
5. **Не обработать ошибки** - задачи в BGTasker будут падать
