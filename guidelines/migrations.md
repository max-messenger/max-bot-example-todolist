# Миграции базы данных

Управление изменениями схемы базы данных с помощью golang-migrate.

## Концепция

TODOLIST использует библиотеку `golang-migrate/migrate/v4` для управления миграциями PostgreSQL. Миграции запускаются автоматически при старте приложения через FX lifecycle hook.

Миграции позволяют:

- Отслеживать изменения схемы БД в виде кода
- Применять изменения последовательно

## Структура файлов миграций

Миграции хранятся в папке `db/migrations`.

### Нейминг миграций

Формат имени: `{timestamp}_{description}.{up|down}.sql`

```bash
20260122182108_create_users_table.up.sql
```

- **timestamp** - `YYYYMMDDHHMMSS` (формат golang-migrate)
- **description** - краткое описание изменения
- **up** - скрипт применения миграции
- **down** - скрипт отката миграции

## Конфигурация

```yaml
migrate:
  pool_name: "upd"  # имя пула для миграций (обычно write pool)
```

Миграции подключаются к приложению через FX:

```go
// internal/app/fx.go
var Modules = fx.Options(
    // ...
    migrate.Module,  // автоматический запуск миграций при старте
    // ...
)
```

## Создание новой миграции

### Через Makefile

```bash
make migrate-new name=create_users_table
```

### Через golang-migrate напрямую

```bash
go tool migrate create -ext sql -dir db/migrations create_users_table
```

Результат будет создано два файла:

```sql
-- db/migrations/20260122182108_create_users_table.up.sql
-- db/migrations/20260122182108_create_users_table.down.sql
```

## Рекомендации по написанию миграций

### 1. IF NOT EXISTS / IF EXISTS

Всегда используйте условные конструкции:

```sql
-- Хорошо
CREATE TABLE IF NOT EXISTS users (...);
DROP TABLE IF EXISTS users;

-- Плохо
CREATE TABLE users (...); -- упадет если таблица существует
DROP TABLE users; -- упадет если таблицы нет
```

### 2. Нейминг колонок и таблиц

Используйте snake_case для названий:

```sql
-- Хорошо
CREATE TABLE user_profiles (
    first_name VARCHAR(255),
    last_name VARCHAR(255)
);

-- Плохо
CREATE TABLE UserProfile (
    FirstName VARCHAR(255),
    LastName VARCHAR(255)
);
```

## Дополнительно

- [golang-migrate documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL documentation](https://www.postgresql.org/docs/)
