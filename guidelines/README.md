# Гайдлайны по разработке

Документация по разработке сервисов с использованием TODOLIST

## Основные документы

- [Архитектура приложения](architecture.md) - обзор принципов и архитектуры
- [Работа с Grace](grace.md) - управление фоновыми сервисами и graceful shutdown
- [PostgreSQL](postgresql.md) - работа с базой данных, read/write реплики, транзакции
- [Redis](redis.md) - работа с Redis
- [Kafka](kafka.md) - продюсеры, консьюмеры и консьюмер-группы
- [BGTasker](bgtasker.md) - асинхронное выполнение задач в фоне
- [Миграции БД](migrations.md) - управление изменениями схемы БД
- [Чек-лист разработки](checklist.md) - пошаговый список действий

## Быстрый старт

Для начала разработки нового сервиса рекомендуется следовать порядку:

1. Ознакомиться с [архитектурой](architecture.md)
2. Изучить примеры:
   - [PostgreSQL](postgresql.md) - для работы с БД
   - [BGTasker](bgtasker.md) - для фоновых задач
   - [Redis](redis.md) - для работы с Redis
3. Следовать [чек-листу](checklist.md) при разработке

## Полезные ссылки

- [Uber FX](https://github.com/uber-go/fx) - Dependency Injection
- [pgx PostgreSQL driver](https://github.com/jackc/pgx) - PostgreSQL драйвер
- [go-redis](https://github.com/redis/go-redis) - Redis клиент
- [franz-go](https://github.com/twmb/franz-go) - Kafka клиент
- [OpenTelemetry Go](https://github.com/open-telemetry/opentelemetry-go) - Observability
