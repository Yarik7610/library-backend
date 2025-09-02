# Library Backend

A modular microservices-based backend for a library-like system with catalog management, user authentication, subscriptions, and notifications.
The project is designed with separation of concerns: each service handles a specific domain, while the API Gateway serves as the single entry point for clients.

Overall project architecture:
<img width="991" height="649" alt="изображение" src="https://github.com/user-attachments/assets/1cdaed8a-470a-4132-9814-f7754e00157d" />

## Core Features

- **API Gateway**

  - Unified access to all backend services
  - Unified access to aggregated Swagger API
  - Auth validation, additional headers incapsultaing

- **User management**

  - Registration and login using JWT
  - User profile information

- **Catalog**

  - Wide range of operations with books (with sorting, ordering, pagination, fetching by categories, case-insensitive pattern matching search)
  - Admin routes for books and author management
  - Redis routes for book views and popularity

- **Subscriptions**

  - Books subscription management

- **Notifications**

  - Email notifications distribution to subscribers (by consuming Kafka's messages)

## Tech Stack

- Docker
- Postgres
- Redis
- Kafka
- Gin, Gorm, Viper, Zap

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Yarik7610/library-backend.git
cd library-backend
```

### 2. Add .env file

Add this file in root of the whole repository:

```env
JWT_SECRET=your_key
MAIL=your_mail # Mail that sends notifications
MAIL_PASSWORD=16_wide_generated_password_in_mail_settings
```

### 3. Start with Docker Compose

To start for the first time (then without `--build` flag):

```bash
docker-compose up --build
```

To stop containers (`-v` flag means with volumes):

```bash
docker-compose up down -v
```

### 4. Create Kafka's topic

```bash
make book-added-topic
```

## API access

Swagger/OpenAPI documentation is included for testing and exploring the endpoints. To test it, visit:

```
http://localhost:80/swagger/index.html
```
