# Redis

Работа с Redis в todolist.

## Концепция

todolist поддерживает работу с несколькими пулами соединений Redis через `Pool`:

```go
type Pool struct {
    cfgs  map[string]*Config   // конфигурации пулов
    pools map[string]*Redis    // активные пулы
}
```

## Конфигурация пулов

Конфигурация пулов определяется в конфигурационном файле:

```yaml
redis:
    default:
      addrs: ["localhost:6379"]
      master_name: ""
      username: ""
      password: ""
      key_prefix: "todolist:"

      max_retries: 3
      min_retry_backoff: 8ms
      max_retry_backoff: 512ms
      max_redirects: 3

      dial_timeout: 5s
      read_timeout: 3s
      write_timeout: 3s

      pool_size: 10
      pool_timeout: 4s
      min_idle_conns: 0
      max_idle_conns: 10
      max_active_conns: 100
      conn_max_idle_time: 5m
      conn_max_lifetime: 30m

      route_randomly: false
      route_by_latency: false
```

### Опции конфигурации

- **addrs** - список адресов Redis серверов
- **master_name** - имя мастера (для Redis Sentinel)
- **key_prefix** - префикс для всех ключей
- **pool_size** - размер пула соединений
- **max_retries** - количество повторных попыток
- **dial_timeout** - таймаут соединения
- **read_timeout** / **write_timeout** - таймауты чтения/записи

## Создание Repository с Redis

```go
package cacherepo

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/max-messenger/max-bot-example-todolist/pkg/connections/rediscli"
)

const (
    RedisPoolName = "default"
)

type CacheRepository struct {
    redis *rediscli.Redis
}

func NewCacheRepository(pool *rediscli.Pool) (*CacheRepository, error) {
    redis, err := pool.GetPool(RedisPoolName)
    if err != nil {
        return nil, fmt.Errorf("get redis pool: %w", err)
    }

    return &CacheRepository{
        redis: redis,
    }, nil
}
```

## Операции с Redis

### Базовые операции

```go
// SET - запись значения с TTL
func (r *CacheRepository) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
    return r.redis.Call(ctx, "set", func(ctx context.Context, ucl redis.UniversalClient, kf rediscli.KeyFormatter) error {
        fullKey := kf.FormatKey(fmt.Sprintf("%s:%s", KeyPrefix, key))

        return ucl.Set(ctx, fullKey, value, ttl).Err()
    })
}

// GET - чтение значения
func (r *CacheRepository) Get(ctx context.Context, key string) (string, error) {
    var result string
    return r.redis.Call(ctx, "get", func(ctx context.Context, ucl redis.UniversalClient, kf rediscli.KeyFormatter) error {
        fullKey := kf.FormatKey(fmt.Sprintf("%s:%s", KeyPrefix, key))

        r, qerr := ucl.Get(ctx, fullKey).Result()
        if qerr != nil {
            if errors.Is(qerr, redis.Nil) {
                return nil
            }
            return qerr
        }

        result = r

       return nil
    })
}

// DEL - удаление ключа
func (r *CacheRepository) Delete(ctx context.Context, key string) error {
    return r.redis.Call(ctx, "delete", func(ctx context.Context, ucl redis.UniversalClient, kf rediscli.KeyFormatter) error {
        fullKey := kf.FormatKey(fmt.Sprintf("%s:%s", KeyPrefix, key))
        return ucl.Del(ctx, fullKey).Err()
    })
}
```

## FX модуль

```go
package cacherepo

import (
    "go.uber.org/fx"
    "go.uber.org/zap"
    "github.com/max-messenger/max-bot-example-todolist/pkg/connections/rediscli"
)

var Module = fx.Module(
    "cacherepo",
    fx.Provide(
        NewCacheRepository,
    ),
    fx.Decorate(func(log *zap.Logger) *zap.Logger {
        return log.Named("cacherepo")
    }),
)
```

## Рекомендации

1. **Key Prefix**: Используйте префикс ключей для избежания конфликтов
2. **TTL**: Всегда устанавливайте TTL для ключей с ограниченным временем жизни
3. **Error Handling**: Обрабатывайте ошибки redis.Nil как отсутствие значения
