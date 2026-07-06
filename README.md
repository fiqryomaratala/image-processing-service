# Image Processing Service

Phase 1 infrastructure for a backend portfolio project built with Go, Gin, RabbitMQ, PostgreSQL, MinIO, Docker, and Docker Compose.

## Current Scope

- Monorepo-ready repository structure
- Backend API and worker entrypoints
- Docker Compose for core infrastructure services
- Placeholder frontend container definition

## Services

- `backend-api`
- `worker`
- `postgres`
- `rabbitmq`
- `minio`

## Quick Start

```bash
docker compose up --build
```

API health check:

```bash
curl http://localhost:8080/health
```
