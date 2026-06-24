# Secure Audit

Распределённая система сбора и анализа событий безопасности (облегчённый SIEM).

## Стек

- **Go 1.25** (Gin, Zap, amqp091-go, golang-jwt/v5)
- **PostgreSQL 16** + JSONB
- **RabbitMQ 3** + плагин management
- **golang-migrate** — автоматические SQL-миграции
- **Grafana**, **Prometheus**
- **Docker Compose**

## Архитектура

```
┌─────────────┐     HTTPS/TLS      ┌───────────┐    AMQP     ┌───────────┐
│   Клиенты   │ ──────────────────> │ Collector │ ──────────> │ RabbitMQ  │
│ (сервисы)   │   API Key + HMAC   │  :8443    │             │  :5672    │
└─────────────┘                     └─────┬─────┘             └─────┬─────┘
                                          │ :9090 /metrics          │
                                     Prometheus               ┌──────▼──────┐
                                                              │  Processor  │
┌─────────────┐                     ┌───────────┐            │ (consumer)  │
│  Grafana    │ ◄────── SQL ──────  │PostgreSQL │ ◄──────────┤  :9091/met. │
│  :3000      │                     │  :5432    │            └──────┬──────┘
└─────────────┘                     └─────┬─────┘                   │
                                          │                  ┌───────┴──────┐
┌─────────────┐        JWT                │                  │  Движок      │
│  Аналитики  │ ──────────────────> ┌─────┴─────┐           │  обнаружения │
│ (браузеры)  │   Bearer token      │    API    │           │ brute_force  │
└─────────────┘                     │  :8081    │           │ suspicious_ip│
                                    └───────────┘           └──────────────┘
```

## Компоненты

| Сервис     | Порт(ы)       | Назначение |
|------------|---------------|------------|
| Collector  | 8443, 9090    | Приём по HTTPS: API Key, HMAC, rate limiting; `/metrics` на 9090 |
| Processor  | 9091          | Очередь → БД, движок обнаружения; `/metrics` на 9091 |
| API        | 8081          | REST API: JWT-аутентификация, журнал аудита запросов |
| PostgreSQL | 5432          | События и оповещения |
| RabbitMQ   | 5672, 15672   | Очередь (веб-интерфейс: 15672) |
| Grafana    | 3000          | Дашборд безопасности |
| pgAdmin    | 5300          | Управление БД |
| Prometheus | 9092          | Метрики (внешний порт; внутри контейнера 9090) |

## Быстрый старт

Единственное требование к машине — установленный **Docker** с Docker Compose. Go, PostgreSQL и RabbitMQ ставить не нужно — всё собирается и запускается в контейнерах.

```bash
# 1. Клонировать репозиторий
git clone https://github.com/Danzhking/secure-audit.git
cd secure-audit

# 2. Секреты
cp .env.example .env
# Отредактируйте .env — замените все CHANGE_ME

# 3. TLS-сертификаты для Collector (одна строка; работает в bash, zsh и PowerShell)
docker run --rm -v ${PWD}/certs:/certs alpine/openssl req -x509 -nodes -newkey rsa:2048 -keyout /certs/server.key -out /certs/server.crt -days 365 -subj "/C=RU/ST=Moscow/O=SecureAudit/CN=collector" -addext "subjectAltName=DNS:collector,DNS:localhost,IP:127.0.0.1"

# 4. Запуск
docker compose up -d --build

# 5. Grafana: http://localhost:3000 (логин/пароль из .env)
```

Миграции БД выполняются **автоматически** при старте Processor.

## Сертификаты

### Разработка (самоподписанный)

Одна строка; работает в bash, zsh и PowerShell:

```bash
docker run --rm -v ${PWD}/certs:/certs alpine/openssl req -x509 -nodes -newkey rsa:2048 -keyout /certs/server.key -out /certs/server.crt -days 365 -subj "/C=RU/ST=Moscow/O=SecureAudit/CN=collector" -addext "subjectAltName=DNS:collector,DNS:localhost,IP:127.0.0.1"
```

### Продакшен

Замените `certs/server.crt` и `certs/server.key` на сертификаты доверенного УЦ. Укажите пути в `.env` (`TLS_CERT`, `TLS_KEY`).

## Проверка работы

### Через curl

```bash
# Получить JWT
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Запрос событий с токеном
curl http://localhost:8081/api/events \
  -H "Authorization: Bearer <токен>"
```

### Демо-пользователи

| Логин    | Пароль   | Роль   |
|----------|----------|--------|
| admin    | admin    | admin  |
| analyst  | analyst  | viewer |
| operator | operator | viewer |

### Postman-коллекция

Файл `docs/secure-audit.postman_collection.json` содержит готовые запросы для ручного тестирования.

**Шаг 1.** Импортируйте коллекцию в Postman (File → Import).  
**Шаг 2.** Откройте настройки коллекции → вкладка *Variables* и задайте значения:

| Переменная   | Значение по умолчанию                          | Откуда брать |
|--------------|------------------------------------------------|--------------|
| `base_url`   | `https://localhost:8443`                       | Collector    |
| `api_url`    | `http://localhost:8081`                        | API-сервис   |
| `api_key`    | `key-1`                                        | `.env` → `COLLECTOR_API_KEYS` |
| `hmac_secret`| `super-secret-hmac-key-change-in-production`   | `.env` → `COLLECTOR_HMAC_SECRET` |
| `jwt_token`  | *(заполняется автоматически)*                  | Тест 7 сохраняет токен при логине |
| `alert_id`   | `1`                                            | ID из ответа на «Список оповещений» |

**Шаг 3.** При самоподписанном сертификате отключите SSL-проверку: *Settings → SSL certificate verification → Off*.  
**Шаг 4.** Выполните **Тест 7** — токен JWT сохранится автоматически для всех последующих запросов.

## Защита Collector

1. **TLS** — шифрование канала (для разработки — самоподписанный сертификат).
2. **API Key** — заголовок `X-API-Key`; сравнение константного времени (защита от timing-атак).
3. **HMAC-SHA256** — заголовок `X-Signature`: HMAC тела запроса в hex.
4. **Rate limiting** — token bucket по IP (10 запр/с, burst 20; неактивные бакеты очищаются автоматически).
5. **Валидация** — теги Gin `binding`.

## Движок обнаружения

На каждое событие применяются правила:

| Правило        | Условие | Важность |
|----------------|---------|----------|
| brute_force    | ≥5 неудачных входов одного пользователя за 10 мин | high |
| suspicious_ip  | ≥3 разных пользователя с одного IP за 5 мин | critical |

Дедупликация оповещений: то же правило не срабатывает повторно в течение 30 минут.

## Модель угроз

### От чего система защищает

- **Перебор паролей** — по порогу неудачных входов на пользователя.
- **Сканирование учётных записей** — один IP бьёт по разным пользователям.
- **Несанкционированный доступ к API** — Collector: ключ + HMAC; чтение данных: JWT.
- **Искажение данных в канале** — TLS на Collector, целостность тела — HMAC.
- **Флуд / DoS** — ограничение частоты на Collector.
- **Несанкционированный просмотр** — аудит запросов к API (пользователь, роль, IP, время).

### Что не покрыто (направления развития)

- **Инсайдеры** — при валидных учётных данных доступ в рамках роли.
- **APT** — только пороговые правила, без ML.
- **Подделка записей в БД** — строки не подписываются отдельно.
- **Компрометация ключей** — ротация без перезапуска не реализована.
- **Сеть Docker** — между сервисами нет mTLS.

## Справочник API

JWT содержит поля `sub` (username) и `role`. Срок действия токена — 24 часа (HS256). Все эндпоинты `/api/*` защищены JWT и логируют запросы через AuditLog middleware (пользователь, роль, IP, метод, статус, латенсия).

### Эндпоинты

| Метод | Путь               | Описание |
|-------|--------------------|----------|
| POST  | /auth/login        | Выдача JWT |
| GET   | /api/events        | Список событий (с фильтрами и пагинацией) |
| GET   | /api/events/:id    | Событие по ID |
| GET   | /api/alerts        | Список оповещений (с фильтрами и пагинацией) |
| PATCH | /api/alerts/:id    | Обновление статуса оповещения |
| GET   | /api/stats         | Агрегированная статистика |

### Фильтры событий

`service`, `event_type`, `severity`, `user_id`, `ip`, `from`, `to`, `page`, `page_size`

### Фильтры оповещений

`rule_name`, `severity`, `status`, `user_id`, `ip`, `from`, `to`, `page`, `page_size`

### Статусы оповещений

`new` → `acknowledged` → `resolved`

## Миграции базы данных

Processor использует [golang-migrate](https://github.com/golang-migrate/migrate) с SQL-файлами, встроенными в бинарник (`//go:embed`). Миграции запускаются автоматически при старте и создают таблицы `security_events` и `alerts` с полным набором индексов.

Каждая миграция имеет шаги `up` и `down`. Файлы миграций:
```
services/processor/internal/repository/migrations/
```

Откат выполняется через CLI golang-migrate, например:
```bash
migrate -path services/processor/internal/repository/migrations \
  -database "$POSTGRES_URL" down 1
```

## Тесты

```bash
cd services/processor && go test ./... -v
cd ../api && go test ./... -v
```
