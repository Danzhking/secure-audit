# Secure Audit

Распределённая система сбора и анализа событий безопасности (облегчённый SIEM).

## Архитектура

```
┌─────────────┐     HTTPS/TLS      ┌───────────┐    AMQP     ┌───────────┐
│   Клиенты   │ ──────────────────> │ Collector │ ──────────> │ RabbitMQ  │
│ (сервисы)   │   API Key + HMAC   │  :8443    │             │  :5672    │
└─────────────┘                     └───────────┘             └─────┬─────┘
                                                                    │
                                                                    ▼
┌─────────────┐                     ┌───────────┐           ┌───────────────┐
│  Grafana    │ ◄────── SQL ──────  │PostgreSQL │ ◄──────── │  Processor    │
│  :3000      │                     │  :5432    │           │ (consumer)    │
└─────────────┘                     └─────┬─────┘           └───────┬───────┘
                                          │                         │
┌─────────────┐        JWT                │                  ┌──────┴───────┐
│  Аналитики  │ ──────────────────> ┌─────┴─────┐           │  Движок      │
│ (браузеры)  │   Bearer token      │    API    │           │  обнаружения │
└─────────────┘                     │  :8081    │           │ brute_force  │
                                    └───────────┘           │ suspicious_ip│
                                                            └──────────────┘
```

## Компоненты

| Сервис     | Порт  | Назначение |
|------------|-------|------------|
| Collector  | 8443  | Приём по HTTPS: API Key, HMAC, ограничение частоты |
| Processor  | —     | Очередь → БД, правила обнаружения |
| API        | 8081  | REST API: JWT, журнал аудита запросов |
| PostgreSQL | 5432  | События и оповещения |
| RabbitMQ   | 5672  | Очередь (веб-интерфейс: 15672) |
| Grafana    | 3000  | Дашборд безопасности |
| pgAdmin    | 5050  | Управление БД |
| Prometheus | 9092  | Метрики (внешний порт; внутри контейнера 9090) |

## Быстрый старт

```bash
# 1. Клонировать репозиторий
git clone https://github.com/Danzhking/secure-audit.git
cd secure-audit

# 2. Секреты
cp .env.example .env
# Отредактируйте .env — замените все CHANGE_ME

# 3. TLS-сертификаты для Collector
docker run --rm -v $(pwd)/certs:/certs alpine/openssl \
  req -x509 -nodes -newkey rsa:2048 \
  -keyout /certs/server.key -out /certs/server.crt \
  -days 365 -subj "/C=RU/ST=Moscow/O=SecureAudit/CN=collector" \
  -addext "subjectAltName=DNS:collector,DNS:localhost,IP:127.0.0.1"

# 4. Запуск
docker compose up -d --build

# 5. Grafana: http://localhost:3000 (логин/пароль из .env)
```

## Аутентификация API

```bash
# Получить JWT
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Запрос событий с токеном
curl http://localhost:8081/api/events \
  -H "Authorization: Bearer <токен>"
```

### Эндпоинты API

| Метод | Путь               | Описание |
|-------|--------------------|----------|
| POST  | /auth/login        | Выдача JWT |
| GET   | /api/events        | Список событий (фильтры) |
| GET   | /api/events/:id    | Событие по ID |
| GET   | /api/alerts        | Список оповещений |
| PATCH | /api/alerts/:id    | Статус оповещения |
| GET   | /api/stats         | Агрегированная статистика |

### Фильтры событий

`service`, `event_type`, `severity`, `user_id`, `ip`, `from`, `to`, `page`, `page_size`

## Защита Collector

1. **TLS** — шифрование канала (для разработки — самоподписанный сертификат).
2. **API Key** — заголовок `X-API-Key` с действующим ключом.
3. **HMAC-SHA256** — заголовок `X-Signature`: HMAC тела запроса в hex.
4. **Rate limiting** — token bucket (10 запр/с, burst 20).
5. **Валидация** — теги Gin `binding` (сообщения валидации от библиотеки могут быть на английском).

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
- **Несанкционированный просмотр** — аудит запросов к API (пользователь, IP, время).

### Что не покрыто (направления развития)

- **Инсайдеры** — при валидных учётных данных доступ в рамках роли.
- **APT** — только пороговые правила, без ML.
- **Подделка записей в БД** — строки не подписываются отдельно.
- **Компрометация ключей** — ротация без перезапуска не реализована.
- **Сеть Docker** — между сервисами нет mTLS.

## Сертификаты

### Разработка (самоподписанный)

```bash
docker run --rm -v $(pwd)/certs:/certs alpine/openssl \
  req -x509 -nodes -newkey rsa:2048 \
  -keyout /certs/server.key -out /certs/server.crt \
  -days 365 -subj "/C=RU/ST=Moscow/O=SecureAudit/CN=collector" \
  -addext "subjectAltName=DNS:collector,DNS:localhost,IP:127.0.0.1"
```

### Продакшен

Замените `certs/server.crt` и `certs/server.key` на сертификаты доверенного УЦ. Укажите пути в `.env` (`TLS_CERT`, `TLS_KEY`).

## Стек

- **Go** (Gin, Zap, amqp091-go)
- **PostgreSQL 16** + JSONB
- **RabbitMQ 3** + плагин management
- **Grafana**, **Prometheus**
- **Docker Compose**

## Тесты

```bash
cd services/processor && go test ./... -v
cd ../api && go test ./... -v
```
