# Secure Audit

Distributed security event collection and analysis system (SIEM-lite).

## Architecture

```
┌─────────────┐     HTTPS/TLS      ┌───────────┐    AMQP     ┌───────────┐
│   Clients   │ ──────────────────> │ Collector │ ──────────> │ RabbitMQ  │
│ (services)  │   API Key + HMAC   │  :8443    │             │  :5672    │
└─────────────┘                     └───────────┘             └─────┬─────┘
                                                                    │
                                                                    ▼
┌─────────────┐                     ┌───────────┐           ┌───────────────┐
│  Grafana    │ ◄────── SQL ──────  │PostgreSQL │ ◄──────── │  Processor    │
│  :3000      │                     │  :5432    │           │ (consumer)    │
└─────────────┘                     └─────┬─────┘           └───────┬───────┘
                                          │                         │
┌─────────────┐        JWT                │                  ┌──────┴───────┐
│  Analysts   │ ──────────────────> ┌─────┴─────┐           │  Detection   │
│ (browsers)  │   Bearer token      │    API    │           │   Engine     │
└─────────────┘                     │  :8081    │           │ brute_force  │
                                    └───────────┘           │ suspicious_ip│
                                                            └──────────────┘
```

## Components

| Service    | Port  | Description                              |
|------------|-------|------------------------------------------|
| Collector  | 8443  | HTTPS ingestion with API Key, HMAC, Rate Limiting |
| Processor  | -     | Consumes queue, stores events, runs detection rules |
| API        | 8081  | REST API with JWT auth, audit logging     |
| PostgreSQL | 5432  | Event and alert storage                   |
| RabbitMQ   | 5672  | Message queue (management UI: 15672)      |
| Grafana    | 3000  | Security dashboard with 10 panels         |
| pgAdmin    | 5050  | Database management UI                    |

## Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/Danzhking/secure-audit.git
cd secure-audit

# 2. Configure secrets
cp .env.example .env
# Edit .env and change all CHANGE_ME values

# 3. Generate TLS certificates
docker run --rm -v $(pwd)/certs:/certs alpine/openssl \
  req -x509 -nodes -newkey rsa:2048 \
  -keyout /certs/server.key -out /certs/server.crt \
  -days 365 -subj "/C=RU/ST=Moscow/O=SecureAudit/CN=collector" \
  -addext "subjectAltName=DNS:collector,DNS:localhost,IP:127.0.0.1"

# 4. Start the system
docker compose up -d

# 5. Open Grafana
# http://localhost:3000 (credentials from .env)
```

## API Authentication

```bash
# Get JWT token
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Use token to query events
curl http://localhost:8081/api/events \
  -H "Authorization: Bearer <token>"
```

### API Endpoints

| Method | Endpoint           | Description                  |
|--------|--------------------|------------------------------|
| POST   | /auth/login        | Get JWT token                |
| GET    | /api/events        | List events (filterable)     |
| GET    | /api/events/:id    | Get event by ID              |
| GET    | /api/alerts        | List alerts (filterable)     |
| PATCH  | /api/alerts/:id    | Update alert status          |
| GET    | /api/stats         | Aggregate statistics         |

### Event Filters

`service`, `event_type`, `severity`, `user_id`, `ip`, `from`, `to`, `page`, `page_size`

## Collector Security Layers

1. **TLS** — All traffic encrypted (self-signed cert for dev, replace for production)
2. **API Key** — Header `X-API-Key` must contain a valid key
3. **HMAC-SHA256** — Header `X-Signature` must contain HMAC of request body
4. **Rate Limiting** — Token Bucket algorithm (10 req/s, burst 20)
5. **Input Validation** — Gin binding tags validate all fields

## Detection Engine

The Processor runs detection rules on every incoming event:

| Rule            | Trigger                                        | Severity |
|-----------------|------------------------------------------------|----------|
| brute_force     | 5+ failed logins by same user in 10 min        | high     |
| suspicious_ip   | 3+ distinct users targeted by same IP in 5 min | critical |

Alerts are deduplicated (same rule won't fire twice within 30 min window).

## Threat Model

### Threats this system protects against

- **Brute force attacks** — Detected by analyzing failed login patterns per user; generates alerts when threshold is exceeded.
- **Credential scanning** — Detected when a single IP targets multiple user accounts in a short time window.
- **Unauthorized API access** — Collector requires API Key + HMAC; API requires JWT Bearer token.
- **Data tampering in transit** — TLS encryption on Collector; HMAC signature verification ensures message integrity.
- **Flood/DoS attacks** — Rate limiting on Collector prevents resource exhaustion.
- **Unauthorized data access** — API audit trail logs every query with user identity, IP, and timestamp.

### Threats NOT covered (future work)

- **Insider threats** — Users with valid credentials can access data within their role scope.
- **Advanced persistent threats** — The detection engine uses simple threshold rules; ML-based anomaly detection is not implemented.
- **Database tampering** — Records in PostgreSQL are not signed; an attacker with DB access could modify records without detection.
- **Key compromise** — API keys and JWT secrets are static; rotation requires service restart.
- **Network-level attacks within Docker** — Inter-service communication inside Docker network is unencrypted (no mTLS).

## Certificate Generation

### Development (self-signed)

```bash
docker run --rm -v $(pwd)/certs:/certs alpine/openssl \
  req -x509 -nodes -newkey rsa:2048 \
  -keyout /certs/server.key -out /certs/server.crt \
  -days 365 -subj "/C=RU/ST=Moscow/O=SecureAudit/CN=collector" \
  -addext "subjectAltName=DNS:collector,DNS:localhost,IP:127.0.0.1"
```

### Production

Replace `certs/server.crt` and `certs/server.key` with certificates from a trusted CA (e.g., Let's Encrypt). Update `TLS_CERT` and `TLS_KEY` in `.env`.

## Tech Stack

- **Go** (Gin, Zap, amqp091-go)
- **PostgreSQL 16** with JSONB metadata
- **RabbitMQ 3** with management plugin
- **Grafana** with PostgreSQL datasource
- **Docker Compose** for orchestration
