# 4-in-a-Row - Local Setup

A real-time multiplayer Connect Four game with Go backend and Next.js frontend.

## Prerequisites

- Go 1.21+
- Node.js 18+
- pnpm 8+
- Docker & Docker Compose
- Make

## Setup Guide

### Part 1: Backend Setup

#### 1. Install Backend Dependencies

```bash
cd backend
make deps
```

#### 2. Choose Your Resource Configuration

**Option A: Local Resources (Recommended for Development)**

Start local PostgreSQL, Redis, Kafka:

```bash
make docker-up
```

**Option B: Cloud Resources**

Skip Docker and use your cloud database/Redis URLs directly.

#### 3. Configure Environment

Create `.env` file in `backend/` directory based on `.env.example`:

```bash
cp .env.example .env
```

Edit `.env` and set:

- `RESOURCE_ENVIRONMENT=local` (for Option A) or `cloud` (for Option B)
- Update database and Redis URLs accordingly
- Set `EVENT_STREAM=redis`

#### 4. Build Backend

```bash
make build
```

#### 5. Run Backend Server

```bash
make run
```

Backend should now be running on http://localhost:8080

#### 6. Run Analytics Consumer (Separate Terminal)

```bash
cd backend
make run-redis-consumer
```

### Part 2: Frontend Setup

#### 1. Install Frontend Dependencies

```bash
cd frontend
pnpm install
```

#### 2. Configure Environment

Create `.env.local` file in `frontend/` directory based on `.env.example`:

```bash
cp .env.example .env.local
```

Update the backend WebSocket/API URLs if needed (defaults should work for local setup).

#### 3. Run Frontend

```bash
pnpm dev
```

Frontend should now be running on http://localhost:3000

## Access Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **WebSocket**: ws://localhost:8080/ws

## API Endpoints

- `GET /leaderboard` - Top 10 players
- `GET /metrics` - Game statistics
- `GET /recent-games` - Last 20 games
- `ws://localhost:8080/ws` - Game WebSocket

## Stopping Services

```bash
# Stop Docker services (if using local resources)
cd backend
make docker-down

# Stop backend/consumer: Ctrl+C
# Stop frontend: Ctrl+C
```

## Development Commands

```bash
# Backend
make help              # Show all commands
make fmt               # Format code
make clean             # Clean build artifacts
make docker-logs       # View Docker logs

# Frontend
pnpm lint              # Lint code
pnpm format            # Format code
pnpm build             # Production build
```

## Troubleshooting

**Services not starting?**

- Ensure Docker is running
- Check ports 5432, 6379, 8080, 3000 are free
- Run `make docker-logs` to view errors

**Connection refused?**

- Verify `.env` configuration matches your resource choice (local/cloud)
- Wait 10 seconds after `make docker-up` for services to initialize

**Build failures?**

- Run `go mod tidy` in backend
- Run `pnpm install` in frontend
