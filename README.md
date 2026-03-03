# Library Backend

A modular microservices-based backend for a library-like system with catalog management, user authentication, subscriptions, and notifications. The project is designed with separation of concerns: each service handles a specific domain, while the API Gateway serves as the single entry point for clients.

Overall project architecture:

<img width="901" height="588" alt="изображение" src="https://github.com/user-attachments/assets/d6473b6f-0ed7-4c8f-8ce2-971c6364de99" />

---

## Core Features

### API Gateway

- Single entry point for all client requests
- Reverse proxy forwarding to downstream microservices
- Aggregated Swagger UI combining docs from all services
- JWT validation and user context propagation via headers

### User Service

- Registration and login with JWT-based authentication
- User profile retrieval
- Admin account seeding on startup

### Catalog Service

- Full CRUD for books and authors (admin-only write operations)
- Advanced book querying: sorting, ordering, pagination, category filtering, case-insensitive search
- Redis-backed book view tracking and popularity ranking

### Subscription Service

- Subscribe and unsubscribe from book categories
- Internal gRPC API to fetch subscribed user emails by category

### Notification Service

- Consumes `book.added` Kafka topic
- Distributes email notifications to all users subscribed to the added book's category
- Worker pool for concurrent email delivery

---

## Tech Stack

- **Docker** — containerization and orchestration via Docker Compose
- **PostgreSQL** — persistent relational storage (one instance per service)
- **Redis** — caching and book popularity tracking (catalog service)
- **Kafka** — async event streaming between catalog and notification services
- **gRPC** — internal synchronous communication between microservices
- **Gin** — HTTP framework
- **GORM** — ORM for PostgreSQL
- **OpenTelemetry** — distributed tracing and metrics instrumentation
- **Jaeger** — tracing backend and UI
- **Prometheus** — metrics collection
- **Grafana** — metrics visualization and dashboards

---

## Architecture Overview

### 3-Layered Architecture

Each microservice follows a strict 3-layer separation:

- **Transport layer** — HTTP handlers and gRPC handlers. Responsible for parsing requests, calling the service layer, and rendering responses. Has no business logic. Works with DTOs.
- **Service layer** — Business logic. Orchestrates repository and external client calls. Knows nothing about HTTP or gRPC. Works with Domains.
- **Repository layer** — Data access. All database queries live here. Works with models.

### gRPC for Internal Communication

Microservice-to-microservice communication uses gRPC instead of HTTP for:

- Strong type safety via Protocol Buffers
- Lower serialization overhead compared to JSON
- Auto-generated client and server code from `.proto` definitions

Proto definitions and generated code are shared via a common module: [library-backend-common](https://github.com/Yarik7610/library-backend-common).

Error mapping is handled in both directions — infrastructure errors are mapped to gRPC status codes on the server side, and gRPC status codes are mapped back to infrastructure errors on the client side. The same is applicapable to HTTP error statuses.

### Observability

**Tracing** is implemented end-to-end using OpenTelemetry. Every incoming HTTP request starts a trace. The trace context is propagated:

- Via HTTP headers when the API Gateway forwards requests to microservices
- Via gRPC metadata when microservices call each other
- Via Kafka message headers when catalog service produces events and notification service consumes them

This means a single `POST /catalog/books` request produces a single trace visible in Jaeger that spans api-gateway → catalog-service → kafka.produce → (async) notification-service (kafka.consume) → subscription-service → user-service.

<img width="1920" height="572" alt="изображение" src="https://github.com/user-attachments/assets/6b94953c-33f0-4cc8-8707-ee93f57d1a82" />

**Metrics** are exposed via Prometheus on each service's `/metrics` endpoint. Grafana dashboards track the 4 Golden Signals:

- **Traffic** — requests per second (RPS)
- **Errors** — 4xx and 5xx error rates
- **Latency** — P50 and P95 response times
- **Saturation** — goroutine count and heap memory usage

<img width="1612" height="805" alt="изображение" src="https://github.com/user-attachments/assets/99646cad-02c7-498d-a9a6-273391c8fdba" />

**Logging** uses structured JSON / Text logging with trace and span IDs injected into every log entry, making it easy to correlate logs with traces in Jaeger.

### Graceful Shutdown

All services handle OS signals (SIGTERM, SIGINT) and shut down gracefully:

1. Stop accepting new HTTP/gRPC requests
2. Wait for in-flight requests to complete
3. Close Kafka readers/writers
4. Close outgoing gRPC client connections
5. Flush and shut down the tracing exporter

### Graceful Bootstrap

Services declare healthchecks in Docker Compose and use `depends_on` with `condition: service_healthy` to ensure dependencies (PostgreSQL, Kafka) are ready before the service starts. Kafka topics are pre-created by a dedicated `kafka-init` container.

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Yarik7610/library-backend.git
cd library-backend
```

### 2. Add `.env` file

Create a `.env` file in the root of the repository:

```env
JWT_SECRET=                          # Secret key for signing JWT tokens
MAIL=                                # Mail address used to send notifications and as seeded admin email
MAIL_PASSWORD=                       # 16-character app password generated in mail settings
ENV=local                            # local or production (default: production)
OTEL_EXPORTER_OTLP_ENDPOINT=         # OTLP endpoint, e.g. http://jaeger:4318
```

> Additional environment variables (ports, DB URLs, Redis config, etc.) are defined per-service in `docker-compose.yml` and can be adjusted there.

### 3. Start

```bash
make up
```

To stop and remove volumes:

```bash
make down
```

---

## API

Swagger UI with aggregated docs from all services is available at:

```
http://localhost:80/swagger/index.html
```

---

## Related

- [library-backend-common](https://github.com/Yarik7610/library-backend-common) — shared proto definitions, generated gRPC code, Kafka topic constants, HTTP route constants, and transport utilities used across all services.
