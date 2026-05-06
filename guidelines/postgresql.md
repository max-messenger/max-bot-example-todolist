# PostgreSQL

Работа с PostgreSQL в todolist, включая read/write реплики и транзакции.

**ВАЖНО!**
`?default_query_exec_mode=exec` необходимо для работы с PostgreSQL в todolist.
Либо `?default_query_exec_mode=simple_protocol` из за особенностей работы с pgbouncer.

## Концепция

todolist поддерживает работу с несколькими пулами соединений PostgreSQL через `Pool`:

```go
type Pool struct {
    cfgs  map[string]*Config   // конфигурации пулов
    pools map[string]*Postgres // активные пулы
}
```

## Конфигурация пулов

Конфигурация пулов определяется в конфигурационном файле:

```yaml
postgres:
    sel:  # read pool (replica)
      database_url: "postgres://user:pass@read-replica-host:5432/db?default_query_exec_mode=exec"
      backoff:
        enabled: true
        max_retries: 3
        initial_interval: 100ms
        max_interval: 1s
        max_elapsed_time: 30s
        multiplier: 2.0

    upd:  # write pool (primary)
      database_url: "postgres://user:pass@write-primary-host:5432/db?default_query_exec_mode=exec"
      backoff:
        enabled: true
        max_retries: 3
```

## Создание Repository с Read/Write репликами

На основе примера `TodoRepository`:

```go
package todorepo

import (
    "context"
    "fmt"

    "github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
    "github.com/max-messenger/max-bot-example-todolist/pkg/connections/postgres"
)

const (
    // Имена пулов из конфигурации
    ReadPGPool  = "sel"  // read replica
    WritePGPool = "upd"  // write primary
)

type Repository struct {
    read  *postgres.Postgres
    write *postgres.Postgres
}

func NewRepository(pool *postgres.Pool) (*Repository, error) {
    // Получаем read pool
    read, err := pool.GetPool(ReadPGPool)
    if err != nil {
        return nil, fmt.Errorf("get pool read: %w", err)
    }

    // Получаем write pool
    write, err := pool.GetPool(WritePGPool)
    if err != nil {
        return nil, fmt.Errorf("get pool write: %w", err)
    }

    return &Repository{
        read:  read,
        write: write,
    }, nil
}
```

## Операции чтения

Используйте `read` пул для SELECT запросов:

```go
import (
    "github.com/huandu/go-sqlbuilder"
    "github.com/max-messenger/max-bot-example-todolist/pkg/connections/postgres"
)

func (r *Repository) List(ctx context.Context) ([]domain.Todo, error) {
    sb := sqlbuilder.Select("id", "user_id", "message", "done", "created_at").
        From("todos")

    result := make([]domain.Todo, 0)
    q, args := sb.Build()

    err := r.read.Call(ctx, "select_todo", func(ctx context.Context, db postgres.Queryable) error {
        rows, err := db.Query(ctx, q, args...)
        if err != nil {
            return err
        }
        defer rows.Close()

        for rows.Next() {
            entity := todoEntity{}
            if err := rows.Scan(
                &entity.ID,
                &entity.UserID,
                &entity.Message,
                &entity.Done,
                &entity.CreatedAt,
            ); err != nil {
                return err
            }
            result = append(result, entity.toDomain())
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    return result, nil
}
```

## Операции записи

Используйте `write` пул для INSERT, UPDATE, DELETE:

```go
func (r *Repository) Add(ctx context.Context, d domain.Todo) error {
    sb := sqlbuilder.InsertInto("todos")
    sb.Cols("user_id", "message")
    sb.Values(d.UserID, d.Message)

    q, args := sb.Build()

    err := r.write.Call(ctx, "insert_todo", func(ctx context.Context, db postgres.Queryable) error {
        _, err := db.Exec(ctx, q, args...)
        return err
    })

    return err
}
```

## Транзакции

Для транзакционных операций используйте `RunInTx`:

```go
import "github.com/jackc/pgx/v5"

func (r *Repository) UpdateWithHistory(ctx context.Context, id int64, newValue string) error {
    txOptions := pgx.TxOptions{
        IsoLevel:   pgx.ReadCommitted,
        AccessMode: pgx.ReadWrite,
    }

    return r.write.RunInTx(ctx, "update_todo_with_history", func(ctx context.Context, tx postgres.Queryable) error {
        // Первая операция
        _, err := tx.Exec(ctx, "UPDATE todos SET message = $1 WHERE id = $2", newValue, id)
        if err != nil {
            return err
        }

        // Вторая операция (в той же транзакции)
        _, err = tx.Exec(ctx, "INSERT INTO todos_history (todo_id, message) VALUES ($1, $2)", id, newValue)
        if err != nil {
            return err
        }

        return nil
    }, txOptions)
}
```

## Пример полного Repository

```go
package todorepo

import (
    "context"
    "fmt"

    "github.com/huandu/go-sqlbuilder"
    "github.com/jackc/pgx/v5"

    "github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
    "github.com/max-messenger/max-bot-example-todolist/pkg/connections/postgres"
)

const (
    ReadPGPool  = "sel"
    WritePGPool = "upd"

    queryNameAddTodo    = "insert_todo"
    queryNameGetTodo    = "get_todo"
    queryNameListTodo   = "list_todos"
    queryNameUpdateTodo = "update_todo"
    queryNameDeleteTodo = "delete_todo"
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

// Create - операция записи
func (r *Repository) Create(ctx context.Context, d domain.Todo) error {
    sb := sqlbuilder.InsertInto("todos")
    sb.Cols("user_id", "message", "done")
    sb.Values(d.UserID, d.Message, d.Done)

    q, args := sb.Build()

    return r.write.Call(ctx, queryNameAddTodo, func(ctx context.Context, db postgres.Queryable) error {
        _, err := db.Exec(ctx, q, args...)
        return err
    })
}

// GetByID - операция чтения
func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.Todo, error) {
    sb := sqlbuilder.Select("id", "user_id", "message", "done", "created_at").
        From("todos").
        Where(sb.Equal("id", id))

    q, args := sb.Build()

    var result *domain.Todo

    err := r.read.Call(ctx, queryNameGetTodo, func(ctx context.Context, db postgres.Queryable) error {
        row := db.QueryRow(ctx, q, args...)
        var entity todoEntity

        if err := row.Scan(
            &entity.ID,
            &entity.UserID,
            &entity.Message,
            &entity.Done,
            &entity.CreatedAt,
        ); err != nil {
            return err
        }

        result = entity.toDomain()
        return nil
    })

    if err != nil {
        return nil, fmt.Errorf("get todo by id: %w", err)
    }

    return result, nil
}

// List - операция чтения
func (r *Repository) List(ctx context.Context) ([]domain.Todo, error) {
    sb := sqlbuilder.Select("id", "user_id", "message", "done", "created_at").
        From("todos")

    q, args := sb.Build()
    result := make([]domain.Todo, 0)

    err := r.read.Call(ctx, queryNameListTodo, func(ctx context.Context, db postgres.Queryable) error {
        rows, err := db.Query(ctx, q, args...)
        if err != nil {
            return err
        }
        defer rows.Close()

        for rows.Next() {
            entity := todoEntity{}
            if err := rows.Scan(
                &entity.ID,
                &entity.UserID,
                &entity.Message,
                &entity.Done,
                &entity.CreatedAt,
            ); err != nil {
                return err
            }
            result = append(result, entity.toDomain())
        }

        return nil
    })

    if err != nil {
        return nil, fmt.Errorf("list todos: %w", err)
    }

    return result, nil
}

// Update - операция записи в транзакции
func (r *Repository) Update(ctx context.Context, id int64, updates map[string]any) error {
    sb := sqlbuilder.Update("todos")
    for field, value := range updates {
        sb.Set(field, value)
    }
    sb.Where(sb.Equal("id", id))

    q, args := sb.Build()

    return r.write.Call(ctx, queryNameUpdateTodo, func(ctx context.Context, db postgres.Queryable) error {
        _, err := db.Exec(ctx, q, args...)
        return err
    })
}

// Delete - операция записи
func (r *Repository) Delete(ctx context.Context, id int64) error {
    sb := sqlbuilder.DeleteFrom("todos")
    sb.Where(sb.Equal("id", id))

    q, args := sb.Build()

    return r.write.Call(ctx, queryNameDeleteTodo, func(ctx context.Context, db postgres.Queryable) error {
        _, err := db.Exec(ctx, q, args...)
        return err
    })
}
```

## FX модуль для Repository

```go
package todorepo

import (
    "go.uber.org/fx"
    "go.uber.org/zap"
    "github.com/max-messenger/max-bot-example-todolist/pkg/connections/postgres"
)

var Module = fx.Module(
    "todorepo",
    fx.Provide(
        NewRepository,
    ),
    fx.Decorate(func(log *zap.Logger) *zap.Logger {
        return log.Named("todorepo")
    }),
)
```

## Рекомендации

1. **Read/Write разделение**: Всегда используйте отдельные пулы для чтения и записи
2. **Именованные запросы**: Используйте осмысленные имена для `queryName` в Call/RunInTx
3. **Транзакции**: Используйте транзакции для операций, требующих атомарности
4. **Обработка ошибок**: Всегда оборачивайте ошибки с контекстом
