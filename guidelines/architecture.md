# Архитектура приложения

Обзор принципов и архитектуры.

## Основные принципы

### Dependency Injection

Использование `uber-go/fx` для управления зависимостями и жизненным циклом компонентов.

```go
import "go.uber.org/fx"

var Module = fx.Module(
    "myservice",
    fx.Provide(
        NewService,
        NewRepository,
    ),
)
```

### Graceful Shutdown

Механизм `Grace` для корректной остановки фоновых сервисов и компонентов.

Подробнее: [Работа с Grace](grace.md)

### Observability

Встроенная телеметрия:

- **OpenTelemetry** - распределенная трассировка
- **Prometheus** - метрики
- **Zap** - структурированное логирование

Все компоненты автоматически обеспечивают сбор метрик и трассировку.

### Connections

Абстракция над подключениями к внешним системам:

- **PostgreSQL** - с поддержкой read/write реплик
- **Redis** - кэширование и хранилище
- **Kafka** - обработка событий

Подробнее: [PostgreSQL](postgresql.md), [Redis](redis.md)

## Компоненты приложения

### Domain Layer

Бизнес-сущности и доменная логика.

```sh
internal/app/domain/
├── todolist.go
└── message.go
```

### Service Layer

Бизнес-логика и оркестрация репозиториев.

```sh
internal/app/services/
├── todolist/
│   ├── todolist.go
│   └── fx.go
└── fx.go
```

### Repository Layer

Работа с базами данных и внешними хранилищами.

```sh
internal/app/repository/
├── todorepo/
│   ├── repository.go
│   ├── entity.go
│   └── fx.go
└── staterepo/
```

### Router Layer

HTTP handlers и API endpoints.

```sh
internal/app/router/
├── todolistctrl/
│   ├── todolistctrl.go
│   └── fx.go
└── router.go
```

## Lifecycle приложения

### Инициализация

1. Загрузка конфигурации
2. Создание FX-приложения с модулями
3. Запуск всех провайдеров
4. Запуск сервисов (Grace ServicePool)

### Работа приложения

- HTTP сервер обрабатывает запросы
- Фоновые сервисы (Grace) выполняют задачи
- Консьюмеры Kafka обрабатывают события

### Завершение работы

1. Сигнал termination (SIGTERM/SIGINT)
2. Graceful shutdown:
   - Остановка HTTP сервера
   - Остановка фоновых сервисов (Grace)
   - Остановка консьюмеров Kafka
   - Закрытие соединений к БД

## Конфигурация

Конфигурация приложения загружается из YAML-файла:

```yaml
# config.yaml
postgres:
  db-sel:
    database_url: "postgres://user:pass@read-host/db"
  db-upd:
    database_url: "postgres://user:pass@write-host/db"

redis:
  main:
    addrs: ["localhost:6379"]

kafka:
  analytic:
    seeds: ["localhost:9092"]

bgtasker:
  tasks:
    analytic:
      async_workers_num: 10
      async_buffer_size: 1000
      runner_type: "drop"
```

## Модульность

todolist построен на модульной архитектуре с использованием `fx.Module`:

```go
// internal/app/fx.go
var Modules = fx.Options(
    // Controllers
    todolistctrl.Module,
    webhookctrl.Module,

    // Services
    services.Module,

    // Repositories
    repository.Module,

    // Connections
    postgres.Module,
    rediscli.Module,
    kafka.Module,

    // migrations
    migrate.Module,

    // Background tasks
    subscriber.Module,
)
```

Каждый модуль - это независимая единица с четкими границами.

## Тестирование

### Unit тесты

Используйте интерфейсы и mocks для тестирования сервисов:

```go
//go:generate go tool mockgen -source=todolist.go -destination=todolist_mocks_test.go -package=todolist_test

type Repository interface {
    Add(ctx context.Context, d domain.Todo) error
    List(ctx context.Context) ([]domain.Todo, error)
}

func TestService_Create(t *testing.T) {
    mockRepo := &MockRepository{}
    service := NewService(mockRepo)

    // тесты
}
```

### Интеграционные тесты

Для тестирования с реальными БД используйте testcontainers.
